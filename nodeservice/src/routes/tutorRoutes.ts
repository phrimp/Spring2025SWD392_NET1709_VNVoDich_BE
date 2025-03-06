import { Router } from "express";
import {
  getTutors,
  getTutor,
  updateTutorProfile,
} from "../controllers/tutorController";
import multer from "multer";

const router = Router();
const upload = multer({ storage: multer.memoryStorage() });

router.get("/", getTutors);
router.get("/:id", getTutor);
router.put("/:id", upload.single("image"), updateTutorProfile);
// router.post("/", createTutor);
// router.delete("/:id", deleteTutor);

export default router;
