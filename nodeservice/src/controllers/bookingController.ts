import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";
import { google } from "googleapis";

import dotenv from "dotenv";
import Stripe from "stripe";
import { isBefore } from "date-fns";

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

  try {
    // Validate request body
    if (
      !courseId ||
      !children_id ||
      !dates ||
      !Array.isArray(dates) ||
      dates.length === 0
    ) {
      res.status(400).json({ message: "Invalid request body" });
      return;
    }

    // Fetch course data more efficiently - only get what's needed
    const course = await prisma.course.findUnique({
      where: {
        id: Number(courseId),
      },
      select: {
        id: true,
        total_lessons: true,
        lessons: {
          select: {
            id: true,
            title: true,
          },
        },
      },
    });

    if (!course) {
      res.status(404).json({ message: "Course not found" });
      return;
    }

    // Generate meet link once before the transaction
    const meetLink = await generateMeetLink();

    const result = await prisma.$transaction(async (tx) => {
      // Create CourseSubscription with bulk schedule creation
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

      // Prepare teaching sessions in bulk instead of one-by-one
      const teachingSessionsData = [];
      let lessonCount = 0;
      const now = new Date(new Date().getTime() + 7 * 60 * 60 * 1000);

      // Tính toán thời gian bắt đầu đã điều chỉnh cho mỗi schedule
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
            status: "Scheduled",
            subscription_id: newBooking.id,
            google_meet_id: meetLink,
            topics_covered: currentLesson?.title || null,
          });

          lessonCount++;
        }
      }

      // Bulk create teaching sessions
      const teachingSessions = await tx.teachingSession.createMany({
        data: teachingSessionsData,
        skipDuplicates: false,
      });

      // Fetch the created sessions
      const createdSessions = await tx.teachingSession.findMany({
        where: {
          subscription_id: newBooking.id,
        },
        orderBy: {
          startTime: "asc",
        },
      });

      return { newBooking, teachingSessions: createdSessions };
    });

    res.json({
      message: "Booking created successfully",
      data: {
        booking: result.newBooking,
        teachingSessions: result.teachingSessions,
      },
    });
  } catch (error) {
    console.error("Error creating booking:", error);
    res.status(500).json({
      message: "Error creating booking and teaching session",
      error: error instanceof Error ? error.message : String(error),
    });
  }
};

const generateMeetLink = async () => {
  const token =
    "ya29.a0AeXRPp43gUm8vSjtFS6DK7A5CFw38S_M2H38hGClVSRuEu3e7HaGJhWVgEJ2GUQYgVSYFJfdlRs0rr4tKdBqKYBGFC4IKc4bWbKHhkQt9hOGaLyB8EymRBMAsKFN4XsPU26INEQPF8pU7z-JLINg6GQdGGOMZKI4pFzVFJ4vgQaCgYKAWsSARISFQHGX2MisCBVKvrCtjfx9MLQDjhRgA0177";

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
      visibility: "private",
      guestsCanSeeOtherGuests: false,
      guestsCanModify: false,
      guestsCanInviteOthers: false,
    },
  });

  const meetLink = meetResponse.data.hangoutLink;

  return meetLink;
};
