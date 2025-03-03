// bookingRoutes.ts
import express from "express";
import {
  createStripePaymentIntent,
  createTrialBooking,
} from "../controllers/bookingController";

const router = express.Router();

router.post("/stripe/payment-intent", createStripePaymentIntent);
router.post("/create-trial-booking", createTrialBooking);

export default router;
