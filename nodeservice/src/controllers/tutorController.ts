import { PrismaClient } from "@prisma/client";
import { Request, Response } from "express";

const prisma = new PrismaClient();

//  Lấy danh sách tất cả tutors với phân trang và bộ lọc
export const getTutors = async (req: Request, res: Response): Promise<void> => {
  try {
    // Lấy query từ request
    const {
      page = 1,
      pageSize = 10,
      isAvailable,
      qualifications,
      teachingStyle,
    } = req.query;

    // Chuyển đổi page và pageSize sang số nguyên
    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    // Tạo bộ lọc dựa trên query parameters
    const filters: any = {};
    if (isAvailable !== undefined)
      filters.is_available = isAvailable === "true";
    if (qualifications)
      filters.qualifications = { contains: qualifications as string };
    if (teachingStyle)
      filters.teaching_style = { contains: teachingStyle as string };

    // Tính toán skip và lấy dữ liệu
    const skip = (pageNum - 1) * pageSizeNum;
    const tutors = await prisma.tutor.findMany({
      where: filters,
      skip,
      take: pageSizeNum,
    });

    // Đếm tổng số bản ghi
    const totalTutors = await prisma.tutor.count({ where: filters });

    // Trả về dữ liệu với phân trang và filter
    res.json({
      message: "Tutors retrieved successfully",
      data: tutors,
      pagination: {
        total: totalTutors,
        page: pageNum,
        pageSize: pageSizeNum,
        totalPages: Math.ceil(totalTutors / pageSizeNum),
      },
    });
  } catch (error: any) {
    res.status(500).json({ message: "Error retrieving tutors", error });
  }
};

//  Lấy thông tin một tutor theo ID
export const getTutor = async (req: Request, res: Response): Promise<void> => {
  const { id } = req.params;

  try {
    const tutor = await prisma.tutor.findUnique({
      where: { id: Number(id) },
    });

    if (!tutor) {
      res.status(404).json({
        message: "Tutor not found",
      });
      return;
    }

    res.json({ message: "Tutor retrieved successfully", data: tutor });
  } catch (error) {
    res.status(500).json({ message: "Error retrieving tutor", error });
  }
};

//  Tạo một tutor mới
export const createTutor = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    let {
      id,
      bio,
      qualifications,
      teaching_style,
      is_available,
      demo_video_url,
      image,
    } = req.body;

    if (!id) {
      id = Math.floor(1 + Math.random() * 9);
    }

    const newTutor = await prisma.tutor.create({
      data: {
        id: Number(id), // ✅ ID là số
        bio,
        qualifications,
        teaching_style,
        is_available: Boolean(is_available),
        demo_video_url: demo_video_url || null,
        image: image || null,
      },
    });

    res.json({ message: "Tutor created successfully", data: newTutor });
  } catch (error: any) {
    console.error("Error creating tutor:", error);
    res.status(500).json({ message: "Error creating tutor", error });
  }
};

//  Xóa một tutor theo ID
export const deleteTutor = async (
  req: Request,
  res: Response
): Promise<void> => {
  const { id } = req.params;

  try {
    const tutor = await prisma.tutor.delete({
      where: { id: Number(id) },
    });

    res.json({ message: "Tutor deleted successfully", data: tutor });
  } catch (error) {
    res.status(500).json({ message: "Error deleting tutor", error });
  }
};
