// bookingRoutes.ts
import express from "express";
import {
  cancelBooking,
  createStripePaymentIntent,
  createTrialBooking,
  getParentBookings,
} from "../controllers/bookingController";
import tutorAuth from "../middleware/tutorAuth";

const router = express.Router();

router.post("/stripe/payment-intent", tutorAuth, createStripePaymentIntent);
router.post("/create-trial-booking", createTrialBooking);
router.get("/parent", tutorAuth, getParentBookings);
router.put("/cancel", tutorAuth, cancelBooking);

export default router;
