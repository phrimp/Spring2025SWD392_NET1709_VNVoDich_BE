// bookingRoutes.ts
import express from "express";
import {
  getTeachingSessions,
  updateTeachingSession,
} from "../controllers/teachingSessionController";

const router = express.Router();

router.get("/", getTeachingSessions);
router.put("/:id", updateTeachingSession);

export default router;
