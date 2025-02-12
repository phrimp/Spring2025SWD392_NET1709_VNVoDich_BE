import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

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
