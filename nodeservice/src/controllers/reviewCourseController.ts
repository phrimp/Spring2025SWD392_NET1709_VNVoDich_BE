import { Request, Response } from "express";
import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

export const addReviewCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { course_id } = req.params;
    const { userId, role, rating, review_content } = req.body;

    const courseId = Number(course_id);
    const parentId = Number(userId);

    if (isNaN(courseId) || isNaN(parentId)) {
      res.status(400).json({ message: "Invalid Course ID or User ID" });
      return;
    }

    // Kiểm tra role của user, chỉ cho phép parent đánh giá
    if (role !== "Parent") {
      res.status(403).json({ message: "Only parents can add course reviews" });
      return;
    }

    // Kiểm tra xem user có ít nhất một children đã đăng ký khóa học không
    const children = await prisma.children.findMany({
      where: { parent_id: parentId },
    });
    if (children.length === 0) {
      res.status(403).json({ message: "User has no registered children" });
      return;
    }

    // Kiểm tra xem có bất kỳ children nào đã đăng ký khóa học này không
    const childrenIds = children.map((child) => child.id);
    const subscription = await prisma.courseSubscription.findFirst({
      where: {
        course_id: courseId,
        children_id: { in: childrenIds },
        status: "Active",
      },
    });

    if (!subscription) {
      res
        .status(403)
        .json({ message: "No active subscription found for user's children" });
      return;
    }

    // Thêm đánh giá khóa học
    const newReview = await prisma.courseReview.create({
      data: {
        course_id: courseId,
        parent_id: parentId,
        rating,
        review_content,
        createAt: new Date(),
      },
    });

    res
      .status(201)
      .json({ message: "Review added successfully", review: newReview });
  } catch (error) {
    console.error("Error adding course review:", error);
    res.status(500).json({ message: "Error adding review", error });
  }
};
