// bookingRoutes.ts
import express from "express";
import { getChildrenTeachingSessions } from "../controllers/teachingSessionController";

const router = express.Router();

router.get("/child/:children_id", getChildrenTeachingSessions);

export default router;
