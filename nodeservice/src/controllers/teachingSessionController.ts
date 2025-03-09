import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

export const getTeachingSessions = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { userId } = req.query;

    const whereClause = userId
      ? {
          subscription: {
            OR: [
              { children_id: Number(userId) },
              { course: { tutor_id: Number(userId) } },
            ],
          },
        }
      : undefined;

    const teachingSessions = await prisma.teachingSession.findMany({
      where: whereClause,
      include: {
        subscription: {
          select: {
            course: {
              include: {
                tutor: {
                  select: {
                    profile: {
                      select: {
                        full_name: true,
                        email: true,
                        phone: true,
                      },
                    },
                  },
                },
              },
            },
          },
        },
      },
    });

    res.json({
      message: "Teaching sessions retrieved successfully",
      data: teachingSessions,
    });
  } catch (error) {
    res
      .status(500)
      .json({ message: "Error retrieving teaching sessions", error });
  }
};

export const rescheduleTeachingSession = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;
    const { startTime, endTime } = req.body;

    const teachingSession = await prisma.teachingSession.findUnique({
      where: { id: Number(id) },
    });

    if (!teachingSession) {
      res.status(404).json({ message: "Teaching session not found" });
      return;
    }

    if (teachingSession.status !== "NotYet") {
      res.status(400).json({ message: "Teaching session has already started" });
      return;
    }

    await prisma.teachingSession.update({
      where: { id: Number(id) },
      data: {
        startTime,
        endTime,
      },
    });

    res.json({
      message: "Teaching sessions rescheduled successfully",
      data: teachingSession,
    });
  } catch (error) {
    res
      .status(500)
      .json({ message: "Error rescheduling teaching sessions", error });
  }
};
