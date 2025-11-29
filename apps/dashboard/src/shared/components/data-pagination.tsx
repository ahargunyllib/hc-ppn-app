import { Button } from "@/shared/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from "lucide-react";

type Props = {
  currentPage: number;
  currentLimit: number;
  setPage: (page: number) => void;
  setLimit: (limit: number) => void;

  totalData: number;
  totalPage: number;
};

export default function DataPagination({
  currentPage,
  currentLimit,
  setPage,
  setLimit,

  totalPage,
}: Props) {
  return (
    <div className="flex items-center justify-end space-x-2 py-4">
      <div className="flex w-full items-center justify-between px-2">
        <div className="flex items-center gap-2">
          <p className="hidden font-medium text-sm sm:block">Rows per page</p>
          <Select
            onValueChange={(value) => {
              setLimit(Number(value));
            }}
            value={`${currentLimit}`}
          >
            <SelectTrigger size="sm">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {[10, 20, 30, 40, 50].map((size) => (
                <SelectItem key={size} value={`${size}`}>
                  {size}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-center space-x-4">
          <div className="flex items-center justify-center font-medium text-sm">
            Page {currentPage} of {totalPage}
          </div>
          <div className="flex items-center space-x-2">
            <Button
              disabled={currentPage === 1}
              onClick={() => setPage(1)}
              size="icon-sm"
              variant="outline"
            >
              <span className="sr-only">Go to first page</span>
              <ChevronsLeft />
            </Button>
            <Button
              disabled={currentPage === 1}
              onClick={() => setPage(currentPage - 1)}
              size="icon-sm"
              variant="outline"
            >
              <span className="sr-only">Go to previous page</span>
              <ChevronLeft />
            </Button>
            <Button
              disabled={currentPage === totalPage}
              onClick={() => setPage(currentPage + 1)}
              size="icon-sm"
              variant="outline"
            >
              <span className="sr-only">Go to next page</span>
              <ChevronRight />
            </Button>
            <Button
              disabled={currentPage === totalPage}
              onClick={() => setPage(totalPage)}
              size="icon-sm"
              variant="outline"
            >
              <span className="sr-only">Go to last page</span>
              <ChevronsRight />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
