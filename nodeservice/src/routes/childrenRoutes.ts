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
router.post("/", tutorAuth, createChild);
router.put("/:id", tutorAuth, updateChild);
router.delete("/:id", tutorAuth, deleteChild);

export default router;
