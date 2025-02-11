import { Router } from "express";
import { getTutorAvailability } from "../controllers/availabilityController";

const router = Router();

router.get("/", getTutorAvailability);

export default router;
