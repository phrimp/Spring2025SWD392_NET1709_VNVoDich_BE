import { Router } from "express";
import {
  getTutorAvailability,
  updateAvailability,
} from "../controllers/availabilityController";
import tutorAuth from "../middleware/tutorAuth";

const router = Router();

router.get("/", tutorAuth, getTutorAvailability);
router.put("/update", tutorAuth, updateAvailability);

export default router;
