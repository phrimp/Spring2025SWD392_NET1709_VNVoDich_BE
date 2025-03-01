import { Request, Response } from "express";
import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

export const getTutorReviewsForParent = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { tutor_id } = req.params;
    const { parent_id } = req.query;

    const tutorId = Number(tutor_id);
    const parentId = parent_id ? Number(parent_id) : null;

    if (isNaN(tutorId)) {
      res.status(400).json({ message: "Invalid Tutor ID" });
      return;
    }

    // Lấy thông tin gia sư từ bảng Tutor
    const tutor = await prisma.tutor.findUnique({
      where: { id: tutorId },
      include: {
        profile: {
          select: { full_name: true },
        },
      },
    });

    if (!tutor) {
      res.status(404).json({ message: "Tutor not found" });
      return;
    }

    // Lấy danh sách đánh giá của gia sư
    const reviews = await prisma.tutorReview.findMany({
      where: {
        tutor_id: tutorId,
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
      tutor_id: tutorId,
      tutor_name: tutor.profile?.full_name || "Unknown Tutor",
      average_rating: Number(averageRating.toFixed(1)), // Trả về dạng số
      reviews: reviews.map((review) => ({
        parent_id: review.parent?.id || null,
        parent_name: review.parent?.profile?.full_name || "Anonymous",
        rating: review.rating,
        review_content: review.review_content,
        createAt: review.createAt, // Nếu trong DB là "createdAt", hãy đổi lại
      })),
    });
  } catch (error) {
    console.error("Error fetching tutor reviews:", error);
    res.status(500).json({ message: "Error fetching reviews", error });
  }
};
