import { Router } from "express";
import {
  getAllUsers,
  getUserById,
  updateUser,
} from "../controllers/manageUserController";
import multer from "multer";
import tutorAuth from "../middleware/tutorAuth";

const router = Router();
const upload = multer({ storage: multer.memoryStorage() });

router.get("/", tutorAuth, getAllUsers);
router.get("/:id", tutorAuth, getUserById);
router.put("/:id", tutorAuth, upload.single("image"), updateUser);

export default router;
