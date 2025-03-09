// bookingRoutes.ts
import express from "express";
import {
  getTeachingSessions,
  rescheduleTeachingSession,
} from "../controllers/teachingSessionController";

const router = express.Router();

router.get("/", getTeachingSessions);
router.put("/reschedule/:id", rescheduleTeachingSession);

export default router;
