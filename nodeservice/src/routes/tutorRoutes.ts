import { Router } from "express";
import { getTutors, getTutor } from "../controllers/tutorController";

const router = Router();

router.get("/", getTutors);
router.get("/:id", getTutor);
// router.post("/", createTutor);
// router.delete("/:id", deleteTutor);

export default router;
