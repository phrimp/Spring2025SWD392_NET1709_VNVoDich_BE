import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

export const getChildrenTeachingSessions = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { children_id } = req.params;

    const teachingSessions = await prisma.teachingSession.findMany({
      where: {
        subscription: {
          children_id: Number(children_id),
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
