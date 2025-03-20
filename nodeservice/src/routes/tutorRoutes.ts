import { Router } from "express";
import {
  getTutors,
  getTutor,
  updateTutorProfile,
  connectTutorAccountToStripe,
  checkConnection,
} from "../controllers/tutorController";
import multer from "multer";
import tutorAuth from "../middleware/tutorAuth";

const router = Router();
const upload = multer({ storage: multer.memoryStorage() });

router.get("/", getTutors);
router.get("/:id", getTutor);
router.put("/:id", upload.single("image"), updateTutorProfile);
router.post(
  "/create-connected-account",
  tutorAuth,
  connectTutorAccountToStripe
);
router.get("/:id/check-connection", checkConnection);
// router.post("/", createTutor);
// router.delete("/:id", deleteTutor);

export default router;
