import { Request, Response } from "express";
import {
  getTutorAvailabilityService,
  updateAvailabilityService,
  getCourseAvailabilityService,
} from "../services/availabilityService";
import { MESSAGES } from "../message/availabilityMessage";

export const getTutorAvailability = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { userId } = req.body;
    const availabilityData = await getTutorAvailabilityService(userId);
    if (!availabilityData) {
      res.json({ message: MESSAGES.tutorNotFound, data: null });
      return;
    }
    res.json({
      message: MESSAGES.availabilityRetrieved,
      data: availabilityData,
    });
  } catch (error) {
    res.status(500).json({
      message: (error as Error).message || MESSAGES.errorRetrievingAvailability,
      error,
    });
  }
};

export const updateAvailability = async (req: Request, res: Response) => {
  try {
    const { userId, ...data } = req.body;
    const updatedAvailability = await updateAvailabilityService(userId, data);
    if (!updatedAvailability) {
      res.json({ message: MESSAGES.tutorNotFound, data: null });
      return;
    }
    res.json({
      message: MESSAGES.availabilityUpdated,
      data: updatedAvailability,
    });
  } catch (error: any) {
    res.status(500).json({
      message: MESSAGES.errorUpdatingAvailability,
      error: (error as Error).message,
    });
  }
};

export const getCourseAvailability = async (req: Request, res: Response) => {
  try {
    const { courseId } = req.params;
    const { type } = req.query;
    const availableDates = await getCourseAvailabilityService(
      Number(courseId),
      type as string
    );
    if (!availableDates) {
      res.status(404).json({ message: MESSAGES.courseNotFound });
      return;
    }
    res.json({
      message: MESSAGES.courseAvailabilityRetrieved,
      data: availableDates,
    });
  } catch (error: any) {
    res.status(500).json({
      message: MESSAGES.errorRetrievingAvailability,
      error: (error as Error).message,
    });
  }
};
