import { Router } from "express";
import {
  getChildren,
  getChild,
  createChild,
  updateChild,
  deleteChild,
} from "../controllers/childrenController";

const router = Router();

router.get("/", getChildren);
router.get("/:id", getChild);
router.post("/", createChild);
router.put("/:id", updateChild);
router.delete("/:id", deleteChild);

export default router;
