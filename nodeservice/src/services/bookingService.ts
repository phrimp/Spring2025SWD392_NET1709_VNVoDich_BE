import { PrismaClient } from "@prisma/client";
import Stripe from "stripe";
import { google } from "googleapis";
import dotenv from "dotenv";
import { isBefore } from "date-fns";
import { BOOKINGMESSAGE } from "../message/bookingMessage";

dotenv.config();
const prisma = new PrismaClient();

if (!process.env.STRIPE_SECRET_KEY) {
  throw new Error(
    "STRIPE_SECRET_KEY is required but was not found in env variables"
  );
}

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY);

export const createStripePaymentIntentService = async (
  amount: number,
  userId: number
) => {
  if (!userId) {
    throw new Error(BOOKINGMESSAGE.USER_ID_REQUIRED);
  }

  const user = await prisma.user.findUnique({
    where: { id: userId },
  });

  if (!user) {
    throw new Error(BOOKINGMESSAGE.USER_NOT_FOUND);
  }

  if (!amount || amount <= 0) {
    amount = 50;
  }

  let customer;
  const doesCustomerExist = await stripe.customers.list({
    email: user.email || `${user.username}@email.com`,
  });

  if (doesCustomerExist.data.length > 0) {
    customer = doesCustomerExist.data[0];
  } else {
    const newCustomer = await stripe.customers.create({
      name: user.username,
      email: user.email || `${user.username}@email.com`,
    });

    customer = newCustomer;
  }

  return await stripe.paymentIntents.create({
    amount: amount * 100,
    currency: "usd",
    automatic_payment_methods: {
      enabled: true,
      allow_redirects: "never",
    },
    customer: customer.id,
    description: "Booking a tutor course",
  });
};

export const createTrialBookingService = async (
  courseId: number,
  children_id: number,
  dates: any[],
  transactionId: string
) => {
  if (
    !courseId ||
    !children_id ||
    !dates ||
    !Array.isArray(dates) ||
    dates.length === 0 ||
    !transactionId
  ) {
    throw new Error(BOOKINGMESSAGE.INVALID_REQUEST_BODY);
  }

  const course = await prisma.course.findUnique({
    where: { id: Number(courseId) },
    select: {
      id: true,
      total_lessons: true,
      price: true,
      lessons: {
        select: { id: true, title: true },
      },
    },
  });

  if (!course) {
    throw new Error(BOOKINGMESSAGE.COURSE_NOT_FOUND);
  }

  const meetLink = await generateMeetLink(children_id);

  return await prisma.$transaction(async (tx) => {
    const newBooking = await tx.courseSubscription.create({
      data: {
        course_id: Number(courseId),
        children_id: Number(children_id),
        status: "Active",
        price: course.price,
        sessions_remaining: course.total_lessons,
        transactionId,
        courseSubscriptionSchedules: {
          createMany: {
            data: dates.map((date) => ({
              startTime: date.startTime,
              endTime: date.endTime,
            })),
          },
        },
      },
      include: {
        courseSubscriptionSchedules: true,
      },
    });

    const totalLessons = course.total_lessons || 1;
    const schedules = newBooking.courseSubscriptionSchedules;
    const weeksNeeded = Math.ceil(totalLessons / schedules.length);

    const teachingSessionsData = [];
    let lessonCount = 0;
    const now = new Date(new Date().getTime() + 7 * 60 * 60 * 1000);

    const adjustedSchedules = schedules.map((schedule) => {
      let adjustedStartTime = new Date(schedule.startTime);
      while (isBefore(adjustedStartTime, now)) {
        adjustedStartTime.setDate(adjustedStartTime.getDate() + 7);
      }
      let adjustedEndTime = new Date(schedule.endTime);
      while (isBefore(adjustedEndTime, now)) {
        adjustedEndTime.setDate(adjustedEndTime.getDate() + 7);
      }
      return { schedule, adjustedStartTime, adjustedEndTime };
    });

    for (
      let week = 0;
      week < weeksNeeded && lessonCount < totalLessons;
      week++
    ) {
      for (const {
        schedule,
        adjustedStartTime,
        adjustedEndTime,
      } of adjustedSchedules) {
        if (lessonCount >= totalLessons) break;

        const currentLesson = course.lessons[lessonCount];

        const startDate = new Date(adjustedStartTime);
        const endDate = new Date(adjustedEndTime);
        startDate.setDate(startDate.getDate() + 7 * week);
        endDate.setDate(endDate.getDate() + 7 * week);

        teachingSessionsData.push({
          startTime: startDate,
          endTime: endDate,
          subscription_id: newBooking.id,
          google_meet_id: meetLink,
          topics_covered: currentLesson?.title || null,
        });

        lessonCount++;
      }
    }

    await tx.teachingSession.createMany({
      data: teachingSessionsData,
      skipDuplicates: false,
    });

    const createdSessions = await tx.teachingSession.findMany({
      where: { subscription_id: newBooking.id },
      orderBy: { startTime: "asc" },
    });

    return { newBooking, teachingSessions: createdSessions };
  });
};

export const getParentBookingsService = async (userId: number) => {
  if (!userId) {
    throw new Error(BOOKINGMESSAGE.USER_ID_REQUIRED);
  }

  return await prisma.courseSubscription.findMany({
    where: {
      children: {
        parent_id: userId,
      },
    },
    include: {
      course: true,
      children: {
        include: {
          profile: {
            select: {
              full_name: true,
            },
          },
        },
      },
    },
  });
};

const generateMeetLink = async (children_id: number) => {
  const children = await prisma.children.findUnique({
    where: { id: children_id },
    select: {
      parent: {
        select: {
          profile: {
            select: {
              googleToken: true,
            },
          },
        },
      },
    },
  });

  const token = children?.parent.profile.googleToken
    ? children?.parent.profile.googleToken
    : process.env.GOOGLE_ACCESS_TOKEN;

  const oauth2Client = new google.auth.OAuth2(
    process.env.GOOGLE_CLIENT_ID,
    process.env.GOOGLE_CLIENT_SECRET,
    process.env.GOOGLE_REDIRECT_URI
  );
  oauth2Client.setCredentials({ access_token: token });

  const calendar = google.calendar({ version: "v3", auth: oauth2Client });

  const meetResponse = await calendar.events.insert({
    calendarId: "primary",
    conferenceDataVersion: 1,
    requestBody: {
      summary: "Tech Talk",
      location: "Google Meet",
      start: { dateTime: new Date().toISOString(), timeZone: "Asia/Kolkata" },
      end: { dateTime: new Date().toISOString(), timeZone: "Asia/Kolkata" },
      attendees: [{ email: "test@example.com" }],
      conferenceData: { createRequest: { requestId: `1-${Date.now()}` } },
    },
  });

  return meetResponse.data.hangoutLink;
};

export const connectTutorAccountToStripeService = async (userId: number) => {
  const user = await prisma.user.findUnique({
    where: { id: userId },
  });

  if (!user || !user.email) {
    throw new Error(BOOKINGMESSAGE.USER_NOT_FOUND);
  }

  const account = await stripe.accounts.create({
    type: "express",
    country: "US",
    email: user.email,
    capabilities: {
      card_payments: { requested: true },
      transfers: { requested: true },
    },
  });

  const destination = account.id;

  await prisma.tutor.update({
    where: { id: userId },
    data: { stripe_account_id: destination },
  });

  const accountLink = await stripe.accountLinks.create({
    account: destination,
    return_url: `${process.env.DOMAIN}/tutor/profile`,
    refresh_url: `${process.env.DOMAIN}/tutor/profile`,
    type: "account_onboarding",
  });

  return { destination, accountLink };
};

export const checkConnectionStatusService = async (userId: number) => {
  console.log(userId);

  if (!userId) {
    throw new Error(BOOKINGMESSAGE.USER_ID_REQUIRED);
  }

  const tutor = await prisma.tutor.findUnique({
    where: { id: userId },
  });

  if (!tutor) {
    throw new Error(BOOKINGMESSAGE.USER_NOT_FOUND);
  }

  if (!tutor.stripe_account_id) {
    return { isConnected: false, description: BOOKINGMESSAGE.NOT_CONNECTED };
  }

  const account = await stripe.accounts.retrieve(tutor.stripe_account_id);
  const isFullyConnected = account.charges_enabled && account.payouts_enabled;

  return {
    isConnected: isFullyConnected,
    description: isFullyConnected
      ? BOOKINGMESSAGE.CONNECTED
      : BOOKINGMESSAGE.NOT_COMPLETE_ONBOARDING,
  };
};
