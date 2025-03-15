import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

export const getTutorsService = async (filters: any, skip: number, pageSizeNum: number) => {
  const tutors = await prisma.tutor.findMany({
    where: filters,
    skip,
    take: pageSizeNum,
    include: {
      profile: {
        select: {
          email: true,
          full_name: true,
          phone: true,
        },
      },
      tutorReviews: true,
    },
  });

  const totalTutors = await prisma.tutor.count({ where: filters });

  return { tutors, totalTutors };
};

export const getTutorService = async (id: number) => {
  return await prisma.tutor.findUnique({
    where: { id },
    include: {
      profile: {
        select: {
          email: true,
          full_name: true,
          phone: true,
          username: true,
        },
      },
      tutorReviews: {
        include: {
          parent: {
            select: {
              profile: {
                select: {
                  full_name: true,
                  picture: true,
                },
              },
            },
          },
        },
      },
    },
  });
};

export const updateTutorProfileService = async (id: number, updateData: any) => {
  return await prisma.tutor.update({
    where: { id },
    data: {
      bio: updateData.bio,
      qualifications: updateData.qualifications,
      teaching_style: updateData.teaching_style,
      demo_video_url: updateData.demo_video_url,
      profile: {
        update: {
          full_name: updateData.full_name,
          phone: updateData.phone,
        },
      },
    },
  });
};
