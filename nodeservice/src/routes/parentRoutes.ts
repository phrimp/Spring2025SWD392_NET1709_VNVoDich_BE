import { Router } from "express";
import {
  getParents,
  getParentById,
  updateParentProfile,
} from "../controllers/parentController";
import multer from "multer";

const router = Router();
const upload = multer({ storage: multer.memoryStorage() });

router.get("/", getParents);
router.get("/:id", getParentById);
router.put("/:id", upload.single("image"), updateParentProfile);

export default router;
