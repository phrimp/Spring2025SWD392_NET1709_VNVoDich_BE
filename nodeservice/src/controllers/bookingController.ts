import { Request, Response } from "express";
import { BOOKINGMESSAGE } from "../message/bookingMessage";
import {
  cancelBookingService,
  createStripePaymentIntentService,
  createTrialBookingService,
  getParentBookingsService,
} from "../services/bookingService";

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
      message: (error as Error).message || BOOKINGMESSAGE.STRIPE_PAYMENT_ERROR,
      error,
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
      message:
        (error as Error).message || BOOKINGMESSAGE.BOOKING_CREATION_ERROR,
      error,
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
      message:
        (error as Error).message || BOOKINGMESSAGE.BOOKING_RETRIEVAL_ERROR,
      error,
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
      message: BOOKINGMESSAGE.BOOKING_IS_CANCELED,
      data: updatedSubscription.updated,
    });
  } catch (error) {

    res.status(500).json({
      message:
        (error as Error).message || BOOKINGMESSAGE.BOOKING_CANCELED_ERROR,

      error,
    });
  }
};
