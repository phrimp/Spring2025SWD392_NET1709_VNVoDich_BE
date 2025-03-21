import { Request, Response } from "express";
import { BOOKINGMESSAGE } from "../message/bookingMessage";
import {
  cancelBookingService,
  createStripePaymentIntentService,
  createTrialBookingService,
  getParentBookingsService,
} from "../services/bookingService";
import { PrismaClient } from "@prisma/client";
import dotenv from "dotenv";
import Stripe from "stripe";

dotenv.config();
const prisma = new PrismaClient();

if (!process.env.STRIPE_SECRET_KEY) {
  throw new Error(
    "STRIPE_SECRET_KEY is required but was not found in env variables"
  );
}

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY);

export const createStripePaymentIntent = async (
  req: Request,
  res: Response
) => {
  try {
    const { amount, userId } = req.body;

    const paymentIntent = await createStripePaymentIntentService(
      Number(amount),
      Number(userId)
    );

    res.json({
      message: "",
      data: { clientSecret: paymentIntent.client_secret },
    });
  } catch (error) {
    res.status(500).json({
      message: BOOKINGMESSAGE.STRIPE_PAYMENT_ERROR,
      error: (error as Error).message,
    });
  }
};

export const createTrialBooking = async (req: Request, res: Response) => {
  try {
    const { courseId, children_id, dates, transactionId } = req.body;
    const result = await createTrialBookingService(
      courseId,
      children_id,
      dates,
      transactionId
    );

    res.json({ message: BOOKINGMESSAGE.BOOKING_SUCCESS, data: result });
  } catch (error) {
    res.status(500).json({
      message: BOOKINGMESSAGE.BOOKING_CREATION_ERROR,
      error: (error as Error).message,
    });
  }
};

export const getParentBookings = async (req: Request, res: Response) => {
  try {
    const { userId } = req.body;
    const bookings = await getParentBookingsService(Number(userId));

    res.json({
      message: BOOKINGMESSAGE.PARENT_BOOKINGS_RETRIEVED,
      data: bookings,
    });
  } catch (error) {
    res.status(500).json({
      message: BOOKINGMESSAGE.BOOKING_RETRIEVAL_ERROR,
      error: (error as Error).message,
    });
  }
};

export const cancelBooking = async (req: Request, res: Response) => {
  try {
    const { subscriptionId, userId } = req.body;

    if (!subscriptionId || !userId) {
      res
        .status(400)
        .json({ message: "Subscription ID or userId is required" });
      return;
    }

    const updatedSubscription = await cancelBookingService(
      subscriptionId,
      userId
    );

    res.json({
      message: "Booking canceled successfully",
      data: updatedSubscription.updated,
    });
  } catch (error) {
    console.error("Error canceling booking:", error);
    res.status(500).json({ message: "Error canceling booking", error });
  }
};
