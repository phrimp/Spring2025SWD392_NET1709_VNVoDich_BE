import { Request, Response } from "express";
import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

export const getCourseReviewsForParent = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { course_id } = req.params;
    const { parent_id } = req.query;

    const courseId = Number(course_id);
    const parentId = parent_id ? Number(parent_id) : null;

    if (isNaN(courseId)) {
      res.status(400).json({ message: "Invalid Course ID" });
      return;
    }

    // Lấy thông tin khóa học và tutor dạy khóa học đó
    const course = await prisma.course.findUnique({
      where: { id: courseId },
      include: {
        tutor: {
          include: {
            profile: {
              select: { full_name: true },
            },
          },
        },
      },
    });

    if (!course) {
      res.status(404).json({ message: "Course not found" });
      return;
    }

    // Lấy danh sách đánh giá của khóa học (có thể lọc theo parent_id nếu có)
    const reviews = await prisma.courseReview.findMany({
      where: {
        course_id: courseId,
        ...(parentId !== null && { parent_id: parentId }),
      },
      include: {
        parent: {
          select: {
            id: true,
            profile: {
              select: { full_name: true },
            },
          },
        },
      },
      orderBy: { createAt: "desc" },
    });

    // Tính điểm trung bình
    const averageRating =
      reviews.length > 0
        ? reviews.reduce((sum, review) => sum + review.rating, 0) /
          reviews.length
        : 0;

    // Trả về kết quả JSON
    res.json({
      course_id: courseId,
      course_title: course.title || "Unknown Course",
      tutor: {
        tutor_id: course.tutor?.id || null,
        tutor_name: course.tutor?.profile?.full_name || "Unknown Tutor",
      },
      average_rating: Number(averageRating.toFixed(1)),
      reviews: reviews.map((review) => ({
        parent_id: review.parent?.id || null,
        parent_name: review.parent?.profile?.full_name || "Anonymous",
        rating: review.rating,
        review_content: review.review_content,
        createAt: review.createAt,
      })),
    });
  } catch (error) {
    console.error("Error fetching course reviews:", error);
    res.status(500).json({ message: "Error fetching reviews", error });
  }
};
