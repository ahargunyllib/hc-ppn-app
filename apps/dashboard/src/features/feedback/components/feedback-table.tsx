import { Button } from "@/shared/components/ui/button";
import {
  PreviewCard,
  PreviewCardPopup,
  PreviewCardTrigger,
} from "@/shared/components/ui/preview-card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/components/ui/table";
import type { Feedback } from "@/shared/types/feedback";
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
} from "@tanstack/react-table";
import {
  BriefcaseIcon,
  CalendarIcon,
  ClockIcon,
  PhoneIcon,
  Star,
  UserIcon,
  UsersIcon,
} from "lucide-react";

type FeedbackTableProps = {
  data: Feedback[];
};

export function FeedbackTable({ data }: FeedbackTableProps) {
  const columns: ColumnDef<Feedback>[] = [
    {
      header: "User",
      cell: ({ row }) => (
        <PreviewCard>
          <PreviewCardTrigger render={<Button size="xs" variant="ghost" />}>
            {row.original.user.phoneNumber}
          </PreviewCardTrigger>
          <PreviewCardPopup align="start">
            <div className="flex w-full flex-col gap-3">
              <div className="flex items-center gap-2">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
                  <UserIcon className="h-5 w-5 text-primary" />
                </div>
                <div className="flex flex-col">
                  <h4 className="font-semibold text-sm">
                    {row.original.user.name}
                  </h4>
                  <div className="flex items-center gap-1 text-muted-foreground text-xs">
                    <PhoneIcon className="h-3 w-3" />
                    <span>{row.original.user.phoneNumber}</span>
                  </div>
                </div>
              </div>

              <div className="h-px w-full bg-border" />

              <div className="flex flex-col gap-2.5">
                {row.original.user.jobTitle && (
                  <div className="flex items-start gap-2">
                    <BriefcaseIcon className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" />
                    <div className="flex flex-col gap-0.5">
                      <span className="text-muted-foreground text-xs">
                        Job Title
                      </span>
                      <span className="font-medium text-sm">
                        {row.original.user.jobTitle}
                      </span>
                    </div>
                  </div>
                )}

                {row.original.user.gender && (
                  <div className="flex items-start gap-2">
                    <UsersIcon className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" />
                    <div className="flex flex-col gap-0.5">
                      <span className="text-muted-foreground text-xs">
                        Gender
                      </span>
                      <span className="font-medium text-sm">
                        {row.original.user.gender}
                      </span>
                    </div>
                  </div>
                )}

                {row.original.user.dateOfBirth && (
                  <div className="flex items-start gap-2">
                    <CalendarIcon className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" />
                    <div className="flex flex-col gap-0.5">
                      <span className="text-muted-foreground text-xs">
                        Date of Birth
                      </span>
                      <span className="font-medium text-sm">
                        {new Date(
                          row.original.user.dateOfBirth
                        ).toLocaleDateString("en-US", {
                          year: "numeric",
                          month: "long",
                          day: "numeric",
                        })}
                      </span>
                    </div>
                  </div>
                )}

                <div className="flex items-start gap-2">
                  <ClockIcon className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" />
                  <div className="flex flex-col gap-0.5">
                    <span className="text-muted-foreground text-xs">
                      Member Since
                    </span>
                    <span className="font-medium text-sm">
                      {new Date(row.original.user.createdAt).toLocaleDateString(
                        "en-US",
                        {
                          year: "numeric",
                          month: "long",
                          day: "numeric",
                        }
                      )}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </PreviewCardPopup>
        </PreviewCard>
      ),
    },
    {
      accessorKey: "rating",
      header: "Rating",
      cell: ({ getValue }) => (
        <div className="flex items-center gap-1">
          <Star className="inline-block h-4 w-4 text-yellow-400" />
          {getValue<number>().toFixed(1)}
        </div>
      ),
    },
    {
      accessorKey: "comment",
      header: "Comment",
      cell: ({ getValue }) => {
        const comment = getValue<string | undefined>();
        return <p className="text-wrap">{comment || "-"}</p>;
      },
    },
    {
      accessorKey: "createdAt",
      header: "Created",
      cell: ({ getValue }) =>
        new Date(getValue<string>()).toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
        }),
    },
  ];

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="overflow-hidden rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                data-state={row.getIsSelected() && "selected"}
                key={row.id}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell className="h-24 text-center" colSpan={columns.length}>
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
