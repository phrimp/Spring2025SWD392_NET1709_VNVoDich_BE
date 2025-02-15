import { Day, DayAvailability, PrismaClient } from "@prisma/client";
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
