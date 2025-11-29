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
import { Star } from "lucide-react";

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
            <div className="flex w-full flex-col gap-2">
              <h4 className="font-medium">Phone Number Details</h4>
              <div className="flex w-full flex-1 flex-col gap-1 rounded-sm border bg-muted p-2 text-muted-foreground text-sm">
                <p>
                  <strong>Label:</strong> {row.original.user.label}
                </p>
                <p>
                  <strong>Assigned To:</strong>{" "}
                  {row.original.user.assignedTo ?? "-"}
                </p>
                <p>
                  <strong>Created At:</strong>{" "}
                  {new Date(row.original.user.createdAt).toLocaleDateString(
                    "en-US",
                    {
                      year: "numeric",
                      month: "short",
                      day: "numeric",
                    }
                  )}
                </p>
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
