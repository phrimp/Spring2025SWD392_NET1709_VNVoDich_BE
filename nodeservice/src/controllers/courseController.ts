import { Request, Response } from "express";
import {
  getCoursesService,
  getCourseByIdService,
  createCourseService,
  updateCourseService,
  deleteCourseService,
  addLessonToCourseService,
  updateLessonService,
  deleteLessonService,
} from "../services/courseService";
import { COURSE_MESSAGES } from "../message/courseMessages";

// Lấy danh sách khóa học
export const getCourses = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const {
      page = 1,
      pageSize = 10,
      subject,
      grade,
      title,
      status,
      userId,
    } = req.query;
    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    const filters: any = {};
    if (subject && subject !== "all")
      filters.subject = { contains: subject as string };
    if (title) filters.title = { contains: title as string };
    if (grade && grade !== "all") filters.grade = parseInt(grade as string, 10);
    if (status) filters.status = status as string;
    if (userId) filters.tutor_id = parseInt(userId as string, 10);

    const skip = (pageNum - 1) * pageSizeNum;
    const { courses, totalCourses } = await getCoursesService(
      filters,
      skip,
      pageSizeNum
    );

    res.json({
      message: COURSE_MESSAGES.COURSES_RETRIEVED,
      data: courses,
      pagination: {
        total: totalCourses,
        page: pageNum,
        pageSize: pageSizeNum,
        totalPages: Math.ceil(totalCourses / pageSizeNum),
      },
    });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_FETCHING,
      error: (error as Error).message,
    });
  }
};

// Lấy chi tiết khóa học
export const getCourse = async (req: Request, res: Response): Promise<void> => {
  try {
    const course = await getCourseByIdService(Number(req.params.id));
    if (!course) {
      res.status(404).json({ message: COURSE_MESSAGES.COURSE_NOT_FOUND });
      return;
    }
    res.json({ message: COURSE_MESSAGES.COURSES_RETRIEVED, data: course });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_FETCHING,
      error: (error as Error).message,
    });
  }
};

// Tạo khóa học
export const createCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    if (!req.body.tutor_id) {
      res.status(400).json({ message: "Tutor Id is required" });
      return;
    }

    const newCourse = await createCourseService(Number(req.body.tutor_id));
    res.json({ message: COURSE_MESSAGES.COURSE_CREATED, data: newCourse });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_CREATING,
      error: (error as Error).message,
    });
  }
};
export const updateCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const courseId = Number(req.params.id);
    const updateData = { ...req.body };

    const updatedCourse = await updateCourseService(courseId, updateData);
    res.json({ message: COURSE_MESSAGES.COURSE_UPDATED, data: updatedCourse });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_UPDATING,
      error: (error as Error).message,
    });
  }
};

//  Xóa khóa học
export const deleteCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const courseId = Number(req.params.id);

    await deleteCourseService(courseId);
    res.json({ message: COURSE_MESSAGES.COURSE_DELETED });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_DELETING,
      error: (error as Error).message,
    });
  }
};

//  Thêm bài học vào khóa học
export const addLessonToCourse = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const courseId = Number(req.params.courseId);
    const lessonData = req.body;

    const updatedCourse = await addLessonToCourseService(courseId, lessonData);
    res.json({ message: COURSE_MESSAGES.LESSON_ADDED, data: updatedCourse });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_CREATING,
      error: (error as Error).message,
    });
  }
};

//  Cập nhật bài học
export const updateLesson = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const courseId = Number(req.params.courseId);
    const lessonId = Number(req.params.lessonId);
    const lessonData = req.body;

    const updatedCourse = await updateLessonService(
      courseId,
      lessonId,
      lessonData
    );
    res.json({ message: COURSE_MESSAGES.LESSON_UPDATED, data: updatedCourse });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_UPDATING,
      error: (error as Error).message,
    });
  }
};

//  Xóa bài học khỏi khóa học
export const deleteLesson = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const courseId = Number(req.params.courseId);
    const lessonId = Number(req.params.lessonId);

    const updatedCourse = await deleteLessonService(courseId, lessonId);
    res.json({ message: COURSE_MESSAGES.LESSON_DELETED, data: updatedCourse });
  } catch (error) {
    res.status(500).json({
      message: COURSE_MESSAGES.ERROR_DELETING,
      error: (error as Error).message,
    });
  }
};
