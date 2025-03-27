import { Request, Response } from "express";
import {
  findTeachingSessions,
  rescheduleTeachingSessionData,
  updateTeachingSessionData,
} from "../services/teachingSessionService";
import { SessionStatus, SessionQuality } from "@prisma/client";
import { TEACHING_SESSION_MESSAGES } from "../message/teachingSessionMessages";

export const getTeachingSessions = async (req: Request, res: Response) => {
  try {
    const { userId } = req.query;

    const teachingSessions = await findTeachingSessions(
      userId ? Number(userId) : undefined
    );

    res.json({
      message: TEACHING_SESSION_MESSAGES.RETRIEVE_SUCCESS,
      data: teachingSessions,
    });
  } catch (error) {
    res.status(500).json({
      message:
        (error as Error).message || TEACHING_SESSION_MESSAGES.RETRIEVE_ERROR,
      error,
    });
  }
};

export const rescheduleTeachingSession = async (
  req: Request,
  res: Response
) => {
  try {
    const { id } = req.params;
    const { startTime, endTime } = req.body;

    const updatedSession = await rescheduleTeachingSessionData(Number(id), {
      startTime,
      endTime,
    });

    if (updatedSession === null) {
      res.status(404).json({ message: TEACHING_SESSION_MESSAGES.NOT_FOUND });
      return;
    }

    res.json({
      message: TEACHING_SESSION_MESSAGES.UPDATE_SUCCESS,
      data: updatedSession,
    });
  } catch (error) {
    console.log(error);

    res.status(500).json({
      message:
        (error as Error).message || TEACHING_SESSION_MESSAGES.UPDATE_ERROR,
      error,
    });
  }
};

export const updateTeachingSession = async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const { status, homework_assigned, rating, teaching_quality, comment } =
      req.body;

    const updatedSession = await updateTeachingSessionData(Number(id), {
      status: status as SessionStatus,
      homework_assigned,
      rating: Number(rating),
      teaching_quality: teaching_quality as SessionQuality,
      comment,
    });

    if (updatedSession === null) {
      res.status(404).json({ message: TEACHING_SESSION_MESSAGES.NOT_FOUND });
      return;
    }

    res.json({
      message: TEACHING_SESSION_MESSAGES.UPDATE_SUCCESS,
      data: updatedSession,
    });
  } catch (error) {
    console.log(error);

    res.status(500).json({
      message:
        (error as Error).message || TEACHING_SESSION_MESSAGES.UPDATE_ERROR,
      error,
    });
  }
};
