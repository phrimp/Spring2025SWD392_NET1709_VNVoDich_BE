import { Day, PrismaClient } from "@prisma/client";
import { addDays, addMinutes, format, isBefore, parseISO, startOfDay } from "date-fns";

const prisma = new PrismaClient();

export const getTutorAvailabilityService = async (userId: number) => {
  const tutor = await prisma.tutor.findUnique({
    where: { id: userId },
    include: { availability: { include: { days: true } } },
  });

  if (!tutor || !tutor.availability) return null;

  const availabilityData: Record<string, any> = { timeGap: tutor.availability.timeGap };
  ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"].forEach((day) => {
    const dayAvailability = tutor.availability?.days.find((d) => d.day === day.toUpperCase());
    availabilityData[day] = {
      isAvailable: !!dayAvailability,
      startTime: dayAvailability ? dayAvailability.startTime.toISOString().slice(11, 16) : "09:00",
      endTime: dayAvailability ? dayAvailability.endTime.toISOString().slice(11, 16) : "17:00",
    };
  });

  return availabilityData;
};

export const updateAvailabilityService = async (userId: number, data: any) => {
  const tutor = await prisma.tutor.findUnique({
    where: { id: userId },
    include: { availability: true },
  });

  if (!tutor) return null;

  const availabilityData = Object.entries(data).flatMap(([day, { isAvailable, startTime, endTime }]: any) => {
    if (isAvailable) {
      const baseDate = new Date().toISOString().split("T")[0];
      return [{ day: day.toUpperCase() as Day, startTime: new Date(`${baseDate}T${startTime}:00Z`), endTime: new Date(`${baseDate}T${endTime}:00Z`) }];
    }
    return [];
  });

  if (tutor.availability) {
    return prisma.availability.update({
      where: { id: tutor.availability.id },
      data: { timeGap: data.timeGap, days: { deleteMany: {}, create: availabilityData } },
    });
  } else {
    return prisma.availability.create({
      data: { tutor_id: tutor.id, timeGap: data.timeGap, days: { create: availabilityData } },
    });
  }
};

export const getCourseAvailabilityService = async (courseId: number, type: string) => {
  const course = await prisma.course.findUnique({
    where: { id: courseId },
    include: {
      tutor: { select: { availability: { select: { days: true, timeGap: true } } } },
      courseSubscriptions: { select: { teachingSessions: { select: { startTime: true, endTime: true } } } },
    },
  });

  if (!course) return null;

  const startDate = startOfDay(new Date().getTime() + 7 * 60 * 60 * 1000);
  const endDate = addDays(startDate, 7);
  const sessions = course.courseSubscriptions.flatMap((sub) => sub.teachingSessions);
  const availableDates = [];

  for (let date = startDate; date <= endDate; date = addDays(date, 1)) {
    const dayOfWeek = format(date, "EEEE").toUpperCase() as Day;
    const dayAvailability = course.tutor.availability?.days.find((d) => d.day === dayOfWeek);
    if (dayAvailability) {
      const dateStr = format(date, "yyyy-MM-dd");
      const slots = generateAvailableTimeSlots({ startTime: dayAvailability.startTime, endTime: dayAvailability.endTime, sessions, dateStr, timeGap: course.tutor.availability?.timeGap, type: type as "Day" | "Week" });
      availableDates.push({ date: dateStr, slots });
    }
  }

  return availableDates;
};

const generateAvailableTimeSlots = ({
  startTime,
  endTime,
  sessions,
  dateStr,
  timeGap = 10,
  duration = 50,
  type = "Week",
}: {
  startTime: Date;
  endTime: Date;
  sessions: { startTime: Date; endTime: Date }[];
  dateStr: string;
  timeGap?: number;
  duration?: number;
  type?: "Day" | "Week";
}) => {
  const slots = [];
  let currentTime = parseISO(`${dateStr}T${startTime.toISOString().slice(11, 16)}:00.000Z`);
  const slotEndTime = parseISO(`${dateStr}T${endTime.toISOString().slice(11, 16)}:00.000Z`);

  if (type === "Day") {
    const now = new Date(new Date().getTime() + 7 * 60 * 60 * 1000);
    if (format(now, "yyyy-MM-dd") === dateStr) {
      currentTime = isBefore(currentTime, now) ? addMinutes(now, timeGap) : currentTime;
    }
  }

  while (currentTime < slotEndTime) {
    const slotEnd = new Date(currentTime.getTime() + duration * 60000 + timeGap * 60000);
    const isSlotAvailable = !sessions.some((booking) => (currentTime >= booking.startTime && currentTime < booking.endTime) || (slotEnd > booking.startTime && slotEnd <= booking.endTime) || (currentTime <= booking.startTime && slotEnd >= booking.endTime));

    if (isSlotAvailable) slots.push(currentTime.toISOString().slice(11, 16));
    currentTime = slotEnd;
  }

  return slots;
};
