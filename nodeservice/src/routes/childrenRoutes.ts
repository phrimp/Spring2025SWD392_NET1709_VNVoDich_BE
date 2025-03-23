import { Router } from "express";
import {
 getChildrenHandler,
 createChildHandler,
 deleteChildHandler,
 getChildHandler,
 updateChildHandler
} from "../controllers/childrenController";
import tutorAuth from "../middleware/tutorAuth";

const router = Router();

router.get("/", tutorAuth, getChildrenHandler);
router.get("/:id", tutorAuth, getChildHandler);

router.post("/", tutorAuth, createChildHandler);
router.put("/:id", tutorAuth, updateChildHandler);
router.delete("/:id", tutorAuth, deleteChildHandler);

export default router;
