import { PrismaClient, SessionStatus, SessionQuality } from "@prisma/client";

const prisma = new PrismaClient();

export const findTeachingSessions = async (userId?: number) => {
  const whereClause = userId
    ? {
        subscription: {
          OR: [
            { children_id: userId },
            { course: { tutor_id: userId } },
          ],
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

  if (!teachingSession) return null;
  if (teachingSession.status !== "NotYet") return "ALREADY_STARTED";

  return await prisma.teachingSession.update({
    where: { id },
    data: {
      ...data,
      rating: data.rating ? Number(data.rating) : undefined,
    },
  });
};
