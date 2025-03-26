import { PrismaClient, SessionStatus, SessionQuality } from "@prisma/client";
import { TEACHING_SESSION_MESSAGES } from "../message/teachingSessionMessages";
import { payoutForTutorService } from "./bookingService";
import { isBefore } from "date-fns";

const prisma = new PrismaClient();

export const findTeachingSessions = async (userId?: number) => {
  const whereClause = userId
    ? {
        subscription: {
          OR: [{ children_id: userId }, { course: { tutor_id: userId } }],
        },
      }
    : undefined;

  return await prisma.teachingSession.findMany({
    where: whereClause,
    include: {
      subscription: {
        select: {
          course: {
            include: {
              tutor: {
                select: {
                  profile: {
                    select: {
                      full_name: true,
                      email: true,
                      phone: true,
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
  });
};

export const updateTeachingSessionData = async (
  id: number,
  data: {
    status?: SessionStatus;
    homework_assigned?: string;
    rating?: number;
    teaching_quality?: SessionQuality;
    comment?: string;
  }
) => {
  const teachingSession = await prisma.teachingSession.findUnique({
    where: { id },
  });

  if (!teachingSession) throw new Error(TEACHING_SESSION_MESSAGES.NOT_FOUND);
  if (teachingSession.status !== "NotYet") {
    throw new Error(TEACHING_SESSION_MESSAGES.ALREADY_STARTED);
  }

  if (new Date(teachingSession.startTime).getDate() !== new Date().getDate()) {
    throw new Error(TEACHING_SESSION_MESSAGES.INVALID_TIME);
  }

  const updatedTeachingSession = await prisma.teachingSession.update({
    where: { id },
    data: {
      ...data,
      rating: data.rating ? Number(data.rating) : undefined,
      subscription: {
        update: {
          sessions_remaining: {
            decrement: 1,
          },
        },
      },
    },
    include: {
      subscription: {
        select: {
          id: true,
          sessions_remaining: true,
          price: true,
          transactionId: true,
          course: {
            select: {
              tutor: {
                select: {
                  stripe_account_id: true,
                },
              },
            },
          },
        },
      },
    },
  });

  const remainingSessions =
    updatedTeachingSession.subscription.sessions_remaining;
  if (remainingSessions === 0) {
    const tutorStripeAccountId =
      updatedTeachingSession.subscription.course.tutor.stripe_account_id;
    const coursePrice = updatedTeachingSession.subscription.price;
    const subscriptionId = updatedTeachingSession.subscription_id;
    // const transactionId = updatedTeachingSession.subscription.transactionId;

    await payoutForTutorService(
      tutorStripeAccountId,
      coursePrice,
      subscriptionId
      // transactionId
    );
  }

  return updatedTeachingSession;
};
export const rescheduleTeachingSessionData = async (
  id: number,
  data: {
    startTime?: string;
    endTime?: string;
  }
) => {
  const teachingSession = await prisma.teachingSession.findUnique({
    where: { id },
  });

  if (!teachingSession) throw new Error(TEACHING_SESSION_MESSAGES.NOT_FOUND);
  if (teachingSession.status !== "NotYet") {
    throw new Error(TEACHING_SESSION_MESSAGES.ALREADY_STARTED);
  }

  const updatedTeachingSession = await prisma.teachingSession.update({
    where: { id },
    data: {
      ...data,
    },
  });

  return updatedTeachingSession;
};
