import { Request, Response } from "express";
import { BOOKINGMESSAGE } from "../message/bookingMessage";
import {
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

    const subscription = await prisma.courseSubscription.findUnique({
      where: { id: Number(subscriptionId) },
      select: {
        id: true,
        status: true,
        transactionId: true,
        course: {
          select: {
            price: true,
          },
        },
        children: {
          select: {
            parent_id: true,
          },
        },
        teachingSessions: {
          select: {
            status: true,
          },
        },
      },
    });

    if (!subscription) {
      res.status(404).json({ message: "Subscription not found" });
      return;
    }

    if (subscription.status === "Canceled") {
      res.status(400).json({ message: "Subscription already canceled" });
      return;
    }

    if (subscription.children.parent_id !== Number(userId)) {
      res.status(403).json({
        message: "You are not authorized to cancel this subscription",
      });
      return;
    }

    // Calculate refund amount
    const totalSessions = subscription.teachingSessions.length;
    const completedSessions = subscription.teachingSessions.filter(
      (session) => session.status !== "NotYet"
    ).length;
    const pricePerSession = subscription.course.price / totalSessions;
    const refundAmount = (totalSessions - completedSessions) * pricePerSession;

    // Xử lý hoàn tiền cho phụ huynh nếu có
    if (refundAmount > 0) {
      if (!subscription.transactionId) {
        throw new Error("Payment Intent ID not found for refund");
      }

      // Tạo refund qua Stripe, hoàn tiền về phương thức thanh toán của phụ huynh
      const refund = await stripe.refunds.create({
        payment_intent: subscription.transactionId,
        amount: Math.round(refundAmount * 100), // Chuyển sang cents
        reason: "requested_by_customer", // Lý do hủy
      });

      await prisma.courseSubscription.update({
        where: { id: Number(subscriptionId) },
        data: {
          refundId: refund.id,
        },
      });
    }

    const updatedSubscription = await prisma.$transaction(async (prisma) => {
      // Update subscription status
      const updated = await prisma.courseSubscription.update({
        where: { id: Number(subscriptionId) },
        data: {
          status: "Canceled",
        },
      });

      // Delete pending sessions
      const deletedSessions = await prisma.teachingSession.deleteMany({
        where: {
          subscription_id: Number(subscriptionId),
          status: "NotYet",
        },
      });

      return { updated, deletedSessionsCount: deletedSessions.count };
    });

    res.json({
      message: "Booking canceled successfully",
      data: updatedSubscription.updated,
    });
  } catch (error) {
    console.error("Error canceling booking:", error);
    res.status(500).json({ message: "Error canceling booking", error });
  }
};
