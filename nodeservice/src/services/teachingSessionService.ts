import { PrismaClient, SessionStatus, SessionQuality } from "@prisma/client";
import { TEACHING_SESSION_MESSAGES } from "../message/teachingSessionMessages";

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
    startTime?: string;
    endTime?: string;
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

  const updatedTeachingSession = await prisma.teachingSession.update({
    where: { id },
    data: {
      ...data,
      rating: data.rating ? Number(data.rating) : undefined,
    },
  });

  if (updatedTeachingSession.status !== "NotYet") {
    await prisma.courseSubscription.update({
      where: { id: updatedTeachingSession.subscription_id },
      data: {
        sessions_remaining: {
          decrement: 1,
        },
      },
    });
  }

  return updatedTeachingSession;
};
