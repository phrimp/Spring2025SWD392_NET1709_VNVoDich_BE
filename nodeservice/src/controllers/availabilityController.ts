import { Request, Response } from "express";
import {
  getTutorAvailabilityService,
  updateAvailabilityService,
  getCourseAvailabilityService,
} from "../services/availabilityService";

export const getTutorAvailability = async (req: Request, res: Response): Promise<void> => {
  try {
    const { userId } = req.body;
    const availabilityData = await getTutorAvailabilityService(userId);
    if (!availabilityData) {
      res.json({ message: "Tutor not found", data: null });
      return;
    }
    res.json({ message: "Availability retrieved successfully", data: availabilityData });
  } catch (error: any) {
    res.status(500).json({ message: "Error retrieving availability", error });
  }
};

export const updateAvailability = async (req: Request, res: Response) => {
  try {
    const { userId, ...data } = req.body;
    const updatedAvailability = await updateAvailabilityService(userId, data);
    if (!updatedAvailability) {
      res.json({ message: "Tutor not found", data: null });
      return;
    }
    res.json({ message: "Availability updated successfully", data: updatedAvailability });
  } catch (error: any) {
    res.status(500).json({ message: "Error updating availability", error });
  }
};

export const getCourseAvailability = async (req: Request, res: Response) => {
  try {
    const { courseId } = req.params;
    const { type } = req.query;
    const availableDates = await getCourseAvailabilityService(Number(courseId), type as string);
    if (!availableDates) {
      res.status(404).json({ message: "Course not found" });
      return;
    }
    res.json({ message: "Course Availability retrieved successfully", data: availableDates });
  } catch (error: any) {
    res.status(500).json({ message: "Error retrieving availability", error });
  }
};
