// bookingController.ts
import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

// Parent booking a course
export const bookCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { parent_id, course_id, children_ids } = req.body;

    if (!parent_id || isNaN(Number(parent_id))) {
      res.status(400).json({ message: "Invalid parent ID" });
      return;
    }

    const parent = await prisma.user.findFirst({
      where: { id: Number(parent_id), role: "Parent" },
    });

    if (!parent) {
      res.status(403).json({ message: "Only parents can book courses." });
      return;
    }

    const course = await prisma.course.findFirst({
      where: { id: Number(course_id), status: "Published" },
    });

    if (!course) {
      res.status(404).json({ message: "Course not found or not available." });
      return;
    }

    const tutor = await prisma.user.findUnique({
      where: { id: course.tutor_id, role: "Tutor" },
    });

    if (!tutor) {
      res.status(404).json({ message: "Tutor not found." });
      return;
    }

    const subscriptions = await Promise.all(
      children_ids.map(async (child_id: number) => {
        return await prisma.courseSubscription.create({
          data: {
            status: "Active",
            sessions_remaining: course.total_lessons,
            course_id: Number(course_id),
            children_id: Number(child_id),
          },
        });
      })
    );

    res.json({
      message: "Course booked successfully",
      tutor,
      subscriptions,
    });
  } catch (error) {
    console.error("Error booking course:", error);
    res.status(500).json({ message: "Error booking course", error });
  }
};

// Get subscriptions for a parent
export const getSubscriptions = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { parent_id } = req.params;

    const children = await prisma.children.findMany({
      where: { parent_id: Number(parent_id) },
      include: {
        courseSubscriptions: {
          include: { course: true },
        },
      },
    });

    res.json({
      message: "Subscriptions retrieved successfully",
      data: children,
    });
  } catch (error) {
    console.error("Error retrieving subscriptions:", error);
    res.status(500).json({ message: "Error retrieving subscriptions", error });
  }
};
