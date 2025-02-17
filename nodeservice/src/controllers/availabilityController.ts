import { Day, PrismaClient } from "@prisma/client";
import { Request, Response } from "express";
import {
  addDays,
  addMinutes,
  format,
  isBefore,
  parseISO,
  startOfDay,
} from "date-fns";

const prisma = new PrismaClient();

export const getTutorAvailability = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { userId } = req.body;

    const tutor = await prisma.tutor.findUnique({
      where: {
        id: userId,
      },
      include: {
        availability: {
          include: { days: true },
        },
      },
    });

    if (!tutor || !tutor.availability) {
      res.json({
        message: "Tutor not found",
        data: null,
      });
      return;
    }

    // Transform the availability data into the format expected by the form
    const availabilityData: Record<string, any> = {
      timeGap: tutor.availability.timeGap,
    };

    [
      "monday",
      "tuesday",
      "wednesday",
      "thursday",
      "friday",
      "saturday",
      "sunday",
    ].forEach((day) => {
      const dayAvailability = tutor.availability?.days.find(
        (d) => d.day === day.toUpperCase()
      );

      availabilityData[day] = {
        isAvailable: !!dayAvailability,
        startTime: dayAvailability
          ? dayAvailability.startTime.toISOString().slice(11, 16)
          : "09:00",
        endTime: dayAvailability
          ? dayAvailability.endTime.toISOString().slice(11, 16)
          : "17:00",
      };
    });

    res.json({
      message: "Availability retrieved successfully",
      data: availabilityData,
    });
  } catch (error: any) {
    res.status(500).json({ message: "Error retrieving availability", error });
  }
};

export const updateAvailability = async (req: Request, res: Response) => {
  try {
    const { userId, ...data } = req.body;

    const tutor = await prisma.tutor.findUnique({
      where: {
        id: userId,
      },
      include: {
        availability: true,
      },
    });

    if (!tutor) {
      res.json({
        message: "Tutor not found",
        data: null,
      });
      return;
    }

    const availabilityData = Object.entries(data).flatMap(
      ([day, { isAvailable, startTime, endTime }]: any) => {
        if (isAvailable) {
          const baseDate = new Date().toISOString().split("T")[0];

          return [
            {
              day: day.toUpperCase() as Day,
              startTime: new Date(`${baseDate}T${startTime}:00Z`),
              endTime: new Date(`${baseDate}T${endTime}:00Z`),
            },
          ];
        }
        return [];
      }
    );

    let updatedAvailability;
    if (tutor.availability) {
      updatedAvailability = await prisma.availability.update({
        where: {
          id: tutor.availability.id,
        },
        data: {
          timeGap: data.timeGap,
          days: {
            deleteMany: {},
            create: availabilityData,
          },
        },
      });
    } else {
      updatedAvailability = await prisma.availability.create({
        data: {
          tutor_id: tutor.id,
          timeGap: data.timeGap,
          days: {
            create: availabilityData,
          },
        },
      });
    }

    res.json({
      message: "Availability updated successfully",
      data: updatedAvailability,
    });
  } catch (error: any) {
    res.status(500).json({ message: "Error retrieving availability", error });
  }
};

export const getCourseAvailability = async (req: Request, res: Response) => {
  try {
    const { courseId } = req.params;

    const course = await prisma.course.findUnique({
      where: {
        id: Number(courseId),
      },
      include: {
        tutor: {
          select: {
            availability: {
              select: {
                days: true,
                timeGap: true,
              },
            },
          },
        },
        courseSubscriptions: {
          select: {
            teachingSessions: {
              select: {
                startTime: true,
                endTime: true,
              },
            },
          },
        },
      },
    });

    if (!course) {
      res.status(404).json({ message: "Course not found" });
      return;
    }

    const startDate = startOfDay(new Date().getTime() + 7 * 60 * 60 * 1000);
    const endDate = addDays(startDate, 7);
    const sessions = course.courseSubscriptions.flatMap(
      (sub) => sub.teachingSessions
    );

    const availableDates = [];

    for (let date = startDate; date <= endDate; date = addDays(date, 1)) {
      const dayOfWeek = format(date, "EEEE").toUpperCase() as Day;
      const dayAvailability = course.tutor.availability?.days.find(
        (d) => d.day === dayOfWeek
      );

      if (dayAvailability) {
        const dateStr = format(date, "yyyy-MM-dd");

        const slots = generateAvailableTimeSlots({
          startTime: dayAvailability.startTime,
          endTime: dayAvailability.endTime,
          sessions,
          dateStr,
          timeGap: course.tutor.availability?.timeGap,
        });

        availableDates.push({
          date: dateStr,
          slots,
        });
      }
    }

    res.json({
      message: "Course Availability retrieved successfully",
      data: { availableDates, course },
    });
  } catch (error: any) {
    res.status(500).json({ message: "Error retrieving availability", error });
  }
};

function generateAvailableTimeSlots({
  startTime,
  endTime,
  sessions,
  dateStr,
  timeGap = 10,
  duration = 50, // Default slot duration in minutes
}: {
  startTime: Date;
  endTime: Date;
  sessions: { startTime: Date; endTime: Date }[];
  dateStr: string;
  timeGap?: number;
  duration?: number;
}) {
  const slots = [];
  let currentTime = parseISO(
    `${dateStr}T${startTime.toISOString().slice(11, 16)}:00.000Z`
  );
  const slotEndTime = parseISO(
    `${dateStr}T${endTime.toISOString().slice(11, 16)}:00.000Z`
  );

  // If the date is today, start from the next available slot after the current time
  const now = new Date(new Date().getTime() + 7 * 60 * 60 * 1000);

  if (format(now, "yyyy-MM-dd") === dateStr) {
    currentTime = isBefore(currentTime, now)
      ? addMinutes(now, timeGap)
      : currentTime;
  }

  while (currentTime < slotEndTime) {
    const slotEnd = new Date(
      currentTime.getTime() + duration * 60000 + timeGap * 60000
    );

    const isSlotAvailable = !sessions.some((booking) => {
      const bookingStart = booking.startTime;
      const bookingEnd = booking.endTime;
      return (
        (currentTime >= bookingStart && currentTime < bookingEnd) ||
        (slotEnd > bookingStart && slotEnd <= bookingEnd) ||
        (currentTime <= bookingStart && slotEnd >= bookingEnd)
      );
    });

    if (isSlotAvailable) {
      slots.push(currentTime.toISOString().slice(11, 16));
    }

    currentTime = slotEnd;
  }

  return slots;
}
