// bookingController.ts
import { CourseSubscription, PrismaClient } from "@prisma/client";
import { Request, Response } from "express";
import { google } from "googleapis";

import dotenv from "dotenv";
import Stripe from "stripe";

const prisma = new PrismaClient();

dotenv.config();

if (!process.env.STRIPE_SECRET_KEY) {
  throw new Error(
    "STRIPE_SECRET_KEY os required but was not found in env variables"
  );
}

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY);

export const createStripePaymentIntent = async (
  req: Request,
  res: Response
): Promise<void> => {
  let { amount } = req.body;

  if (!amount || amount <= 0) {
    amount = 50;
  }

  try {
    const paymentIntent = await stripe.paymentIntents.create({
      amount,
      currency: "usd",
      automatic_payment_methods: {
        enabled: true,
        allow_redirects: "never",
      },
    });

    res.json({
      message: "",
      data: {
        clientSecret: paymentIntent.client_secret,
      },
    });
  } catch (error) {
    res
      .status(500)
      .json({ message: "Error creating stripe payment intent", error });
  }
};

export const createTrialBooking = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { courseId, children_id, dates } = req.body;

  console.log(dates);

  try {
    if (!courseId || !children_id || !dates || !Array.isArray(dates)) {
      res.status(400).json({ message: "Invalid request body" });
      return;
    }

    const course = await prisma.course.findUnique({
      where: {
        id: Number(courseId),
      },
      include: {
        lessons: true,
      },
    });

    if (!course) {
      res.status(404).json({ message: "Course not found" });
      return;
    }

    const result = await prisma.$transaction(async (tx) => {
      // Create CourseSubscription
      const newBooking = await tx.courseSubscription.create({
        data: {
          course_id: Number(courseId),
          children_id: Number(children_id),
          status: "Active",
          sessions_remaining: course?.total_lessons,
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

      const teachingSessions = [];

      const totalLessons = course?.total_lessons || 1;

      const schedules = newBooking.courseSubscriptionSchedules;

      const weeksNeeded = Math.ceil(totalLessons / schedules.length);

      let lessonCount = 0;

      const meetLink = await generateMeetLink();

      for (let week = 0; week < weeksNeeded; week++) {
        // For each schedule in a week
        for (const schedule of schedules) {
          // Stop if we've created enough sessions
          if (lessonCount >= totalLessons) {
            break;
          }

          const currentLesson = course.lessons[lessonCount];
          // Calculate the date for this week (add 7 days for each week)
          const startDate = new Date(schedule.startTime);
          const endDate = new Date(schedule.endTime);

          // Add weeks (7 * week days)
          startDate.setDate(startDate.getDate() + 7 * week);
          endDate.setDate(endDate.getDate() + 7 * week);

          // Create teaching session for this date
          const teachingSession = await tx.teachingSession.create({
            data: {
              startTime: startDate,
              endTime: endDate,
              status: "Scheduled", // Default status for new teaching sessions
              subscription_id: newBooking.id,
              google_meet_id: meetLink,
              topics_covered: currentLesson?.title,
            },
          });

          teachingSessions.push(teachingSession);
          lessonCount++;
        }
      }

      return { newBooking, teachingSessions };
    });

    res.json({
      message: "Booking created successfully",
      data: {
        booking: result.newBooking,
        teachingSessions: result.teachingSessions,
      },
    });
  } catch (error) {
    res
      .status(500)
      .json({ message: "Error creating booking and teaching session", error });
  }
};

const generateMeetLink = async () => {
  const token =
    "ya29.a0AeXRPp5Jhfi-N1_S15aPPXMYzBPAp1mWxoFAncWQWoqBMoNHJIG2Bxh6K_GsD1TMF7Xy_dJse_58oEqJo_dxJwx_cTuSzIdKibeKarjULa631ilSTAmKXSIvTy7NM_NN1xeYKkPUhcqYZpcRDkcf459AQGtcQj349oP00C7FaCgYKATASARISFQHGX2MiFoJuRm62-lnbmmFm_8Db1Q0175";

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
      summary: "Tech Talk with Arindam",
      location: "Google Meet",
      description: "Demo event for Arindam's Blog Post.",
      start: {
        dateTime: "2024-03-14T19:30:00+05:30",
        timeZone: "Asia/Kolkata",
      },
      end: {
        dateTime: "2024-03-14T20:30:00+05:30",
        timeZone: "Asia/Kolkata",
      },
      attendees: [{ email: "quansieuquay2013@gmail.com" }],
      conferenceData: {
        createRequest: { requestId: `1-${Date.now()}` },
      },
    },
  });

  const meetLink = meetResponse.data.hangoutLink;

  return meetLink;
};
