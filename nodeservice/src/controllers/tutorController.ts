import { Request, Response } from "express";
import { getTutorsService, getTutorService, updateTutorProfileService } from "../services/tutorService";
import { tutorMessages } from "../message/tutorMessage";

// Lấy danh sách tutors
export const getTutors = async (req: Request, res: Response): Promise<void> => {
  try {
    const {
      page = 1,
      pageSize = 10,
      qualifications,
      teachingStyle,
      isAvailable,
      email,
      fullName,
      phone,
    } = req.query;

    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    const filters: any = {};
    if (qualifications) filters.qualifications = { contains: qualifications as string };
    if (teachingStyle) filters.teaching_style = { contains: teachingStyle as string };
    if (isAvailable !== undefined) filters.is_available = isAvailable === "true";
    if (email) filters.user = { email: { contains: email as string } };
    if (fullName) filters.user = { full_name: { contains: fullName as string } };
    if (phone) filters.user = { phone: { contains: phone as string } };

    const skip = (pageNum - 1) * pageSizeNum;
    const { tutors, totalTutors } = await getTutorsService(filters, skip, pageSizeNum);

    res.json({
      message: tutorMessages.SUCCESS.GET_TUTORS,
      data: tutors,
      pagination: {
        total: totalTutors,
        page: pageNum,
        pageSize: pageSizeNum,
        totalPages: Math.ceil(totalTutors / pageSizeNum),
      },
    });
  } catch (error) {
    console.error(tutorMessages.ERROR.GET_TUTORS, error);
    res.status(500).json({ message: tutorMessages.ERROR.GET_TUTORS, error });
  }
};

// Lấy thông tin một tutor
export const getTutor = async (req: Request, res: Response): Promise<void> => {
  const { id } = req.params;

  try {
    const tutor = await getTutorService(Number(id));

    if (!tutor) {
      res.status(404).json({ message: tutorMessages.ERROR.TUTOR_NOT_FOUND });
      return;
    }

    res.json({ message: tutorMessages.SUCCESS.GET_TUTOR, data: tutor });
  } catch (error) {
    console.error(tutorMessages.ERROR.GET_TUTOR, error);
    res.status(500).json({ message: tutorMessages.ERROR.GET_TUTOR, error });
  }
};

// Cập nhật thông tin tutor
export const updateTutorProfile = async (req: Request, res: Response): Promise<void> => {
  const { id } = req.params;
  const updateData = { ...req.body };

  try {
    const tutor = await updateTutorProfileService(Number(id), updateData);

    if (!tutor) {
      res.status(404).json({ message: tutorMessages.ERROR.TUTOR_NOT_FOUND });
      return;
    }

    res.json({ message: tutorMessages.SUCCESS.UPDATE_TUTOR, data: tutor });
  } catch (error) {
    console.error(tutorMessages.ERROR.UPDATE_TUTOR, error);
    res.status(500).json({ message: tutorMessages.ERROR.UPDATE_TUTOR, error });
  }
};
