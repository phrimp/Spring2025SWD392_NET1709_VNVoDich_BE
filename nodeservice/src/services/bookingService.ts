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

export const createStripePaymentIntentService = async (amount: number) => {
  if (!amount || amount <= 0) {
    amount = 50;
  }

  return await stripe.paymentIntents.create({
    amount,
    currency: "usd",
    automatic_payment_methods: {
      enabled: true,
      allow_redirects: "never",
    },
  });
};

export const createTrialBookingService = async (
  courseId: number,
  children_id: number,
  dates: any[]
) => {
  if (
    !courseId ||
    !children_id ||
    !dates ||
    !Array.isArray(dates) ||
    dates.length === 0
  ) {
    throw new Error(BOOKINGMESSAGE.INVALID_REQUEST_BODY);
  }

  const course = await prisma.course.findUnique({
    where: { id: Number(courseId) },
    select: {
      id: true,
      total_lessons: true,
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
        sessions_remaining: course.total_lessons,
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
