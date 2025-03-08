// bookingRoutes.ts
import express from "express";
import { getTeachingSessions } from "../controllers/teachingSessionController";

const router = express.Router();

router.get("/", getTeachingSessions);

export default router;
