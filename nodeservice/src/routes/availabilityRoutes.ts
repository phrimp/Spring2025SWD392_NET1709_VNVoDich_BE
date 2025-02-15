import { Router } from "express";
import { getTutorAvailability } from "../controllers/availabilityController";
import tutorAuth from "../middleware/tutorAuth";

const router = Router();

router.get("/", tutorAuth, getTutorAvailability);

export default router;
