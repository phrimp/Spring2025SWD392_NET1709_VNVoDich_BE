import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

export const getTeachingSessions = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { userId } = req.query;

    const teachingSessions = userId
      ? await prisma.teachingSession.findMany({
          include: {
            subscription: {
              select: {
                course: true,
              },
            },
          },
        })
      : await prisma.teachingSession.findMany({
          where: {
            subscription: {
              OR: [
                { children_id: Number(userId) },
                {
                  course: {
                    tutor_id: Number(userId),
                  },
                },
              ],
            },
          },
          include: {
            subscription: {
              select: {
                course: true,
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
