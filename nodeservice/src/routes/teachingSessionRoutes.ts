// bookingRoutes.ts
import express from "express";
import {
  getTeachingSessions,
  rescheduleTeachingSession,
  updateTeachingSession,
} from "../controllers/teachingSessionController";

const router = express.Router();

router.get("/", getTeachingSessions);
router.put("/:id", updateTeachingSession);
router.put("/:id/reschedule", rescheduleTeachingSession);

export default router;
