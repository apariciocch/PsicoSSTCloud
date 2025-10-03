import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';

@Injectable()
export class CompaniesService {
  constructor(private prisma: PrismaService) {}

  async findAll() {
    return this.prisma.company.findMany();
  }

  async create(data: any) {
    return this.prisma.company.create({ data });
  }
}
