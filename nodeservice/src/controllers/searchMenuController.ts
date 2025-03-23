import { Request, Response } from "express";
import {
  findTutors,
  findCourses,
  filterCoursesByPrice,
  filterCoursesByGrade,
  filterCoursesBySubject,
  filterTutorsByRating,
} from "../services/searchService";
import { SEARCH_MESSAGES } from "../message/searchMessages";

export const searchTutors = async (req: Request, res: Response) => {
  try {
    const { query, page = 1, pageSize = 10 } = req.query;

    const tutors = await findTutors(
      query as string,
      Number(page),
      Number(pageSize)
    );
    res.json({ message: SEARCH_MESSAGES.TUTORS_SUCCESS, data: { tutors } });
  } catch (error) {
    res.status(500).json({
      message: (error as Error).message || SEARCH_MESSAGES.ERROR_SEARCH_TUTORS,
      error,
    });
  }
};

export const searchCourses = async (req: Request, res: Response) => {
  try {
    const { query, page = 1, pageSize = 10 } = req.query;

    const courses = await findCourses(
      query as string,
      Number(page),
      Number(pageSize)
    );
    res.json({ message: SEARCH_MESSAGES.COURSES_SUCCESS, data: { courses } });
  } catch (error) {
    res.status(500).json({
      message: (error as Error).message || SEARCH_MESSAGES.ERROR_SEARCH_COURSES,
      error,
    });
  }
};

export const filterByPrice = async (req: Request, res: Response) => {
  try {
    const { minPrice, maxPrice } = req.query;

    const courses = await filterCoursesByPrice(
      Number(minPrice),
      Number(maxPrice)
    );
    res.json({ message: SEARCH_MESSAGES.FILTER_PRICE_SUCCESS, data: courses });
  } catch (error) {
    res.status(500).json({
      message: (error as Error).message || SEARCH_MESSAGES.ERROR_FILTER_PRICE,
      error,
    });
  }
};

export const searchByGrade = async (req: Request, res: Response) => {
  try {
    const { grade } = req.query;

    const courses = await filterCoursesByGrade(Number(grade));
    res.json({ message: SEARCH_MESSAGES.FILTER_GRADE_SUCCESS, data: courses });
  } catch (error) {
    res.status(500).json({
      message: (error as Error).message || SEARCH_MESSAGES.ERROR_FILTER_GRADE,
      error,
    });
  }
};

export const searchBySubject = async (req: Request, res: Response) => {
  try {
    const { subject } = req.query;

    const courses = await filterCoursesBySubject(subject as string);
    res.json({
      message: SEARCH_MESSAGES.FILTER_SUBJECT_SUCCESS,
      data: courses,
    });
  } catch (error) {
    res.status(500).json({
      message: (error as Error).message || SEARCH_MESSAGES.ERROR_FILTER_SUBJECT,
      error,
    });
  }
};

export const filterTutorsByRatings = async (req: Request, res: Response) => {
  try {
    const minRating = Number(req.query.minRating) || 0;

    const tutorsWithAvgRating = await filterTutorsByRating(minRating);
    res.json({
      message: SEARCH_MESSAGES.FILTER_TUTORS_RATING_SUCCESS,
      data: tutorsWithAvgRating,
    });
  } catch (error) {
    res.status(500).json({
      message: (error as Error).message || SEARCH_MESSAGES.ERROR_FILTER_TUTORS,
      error,
    });
  }
};
