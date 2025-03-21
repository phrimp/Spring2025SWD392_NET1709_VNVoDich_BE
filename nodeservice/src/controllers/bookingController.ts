import { Request, Response } from "express";
import { BOOKINGMESSAGE } from "../message/bookingMessage";
import {
  createStripePaymentIntentService,
  createTrialBookingService,
  getParentBookingsService,
} from "../services/bookingService";

export const createStripePaymentIntent = async (
  req: Request,
  res: Response
) => {
  try {
    const { amount } = req.body;
    const paymentIntent = await createStripePaymentIntentService(amount);

    res.json({
      message: BOOKINGMESSAGE.PAYMENT_SUCCESS,
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
    const { courseId, children_id, dates } = req.body;
    const result = await createTrialBookingService(
      courseId,
      children_id,
      dates
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
