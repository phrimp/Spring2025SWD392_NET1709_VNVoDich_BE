import { Router } from "express";
import {
  getCourseAvailability,
  getTutorAvailability,
  updateAvailability,
} from "../controllers/availabilityController";
import tutorAuth from "../middleware/tutorAuth";

const router = Router();

router.get("/", tutorAuth, getTutorAvailability);
router.get("/course/:courseId", getCourseAvailability);
router.put("/update", tutorAuth, updateAvailability);

export default router;
