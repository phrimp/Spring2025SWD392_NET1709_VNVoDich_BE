import { Router } from "express";
import {
  handleGetParentById,
  handleGetParents,
  handleUpdateParentProfile,
} from "../controllers/parentController";
import multer from "multer";

const router = Router();
const upload = multer({ storage: multer.memoryStorage() });

router.get("/", handleGetParents);
router.get("/:id", handleGetParentById);
router.put("/:id", upload.single("image"), handleUpdateParentProfile);

export default router;
