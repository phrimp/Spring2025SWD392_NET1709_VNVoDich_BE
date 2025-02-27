import { Router } from "express";
import {
  getChildren,
  getChild,
  createChild,
  updateChild,
  deleteChild,
} from "../controllers/childrenController";
import tutorAuth from "../middleware/tutorAuth";

const router = Router();

router.get("/", tutorAuth, getChildren);
router.get("/:id", tutorAuth, getChild);
router.post("/", createChild);
router.put("/:id", updateChild);
router.delete("/:id", deleteChild);

export default router;
