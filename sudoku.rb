#!/usr/bin/env ruby

# Solves all 120 puzzles in Wayne Gould's Extreme Su Doku 2
# See http://www.sudokuwiki.org for explanation of strategies
# Note: I do not use hidden strategies since these seem to be mirror
# images of longer naked strategies and I personally am more likely
# to spot a long naked sequence than a short hidden one.

class Cell

    attr_reader :board, :row, :col, :possibles, :ref
    
    def initialize(board, row, col, ref)
        @board = board
        @row = row
        @col = col
        @ref = ref
        @possibles = '123456789'
    end
    
    def box()
        (@row / 3) * 3 + (@col / 3)
    end
    
    def has?(possible)
        @possibles.length > 1 and @possibles.include? possible
    end
    
    def solve(value)
        @possibles = value
    end
    
    def remove(value)
        found = false
        if @possibles.length > 1
            found = @possibles.sub!(value, '')
            if @possibles.length == 1
                @board.solve(self, @possibles)
            end
        end
        found
    end
    
    def remove_possibles(values)
        found = false        
        values.each_char { |value| found = true if remove(value) }
        found
    end
    
    def to_s()
        "Cell #{@ref} #{@row}:#{@col}=#{@possibles}"
    end
    
end 

class Board

    def initialize(data)
        @cells = []
        i = 0
        (0..8).each { |row| (0..8).each { |col| @cells << Cell.new(self, row, col, i); i += 1 }}
        i = 0
        (0..8).each { |row| (0..8).each { |col| 
            value = data[i,1]
            solve(cell(row, col), value) if '123456789'.include? value 
            i += 1
        }}
            
    end
    
    def to_s()
        rows = []
        (0..8).each do |row|
            rows << cells_by_row(row).each.map { |cell| cell.possibles }
        end
        (0..8).each do |col|
            max = 0
            (0..8).each do |row|
                max = rows[row][col].length if rows[row][col].length > max
            end
            (0..8).each do |row|
                rows[row][col] = rows[row][col] + " " * (max - rows[row][col].length)
            end
        end
        s = []
        s << '=' * 50
        i = 0
        rows.each do |row|
            i += 1
            s << "#{row.join("|")}"
            s << '-' * 50 if i == 3 or i == 6
        end
        s << '=' * 50
        s.join("\n")
    end
    
    def solve(solved_cell, value)
        solved_cell.solve(value)
        @cells.each { |cell| cell.remove(value) if cell.row == solved_cell.row or cell.col == solved_cell.col or cell.box == solved_cell.box }
    end
    
    def cell(row, col)
        @cells[row * 9 + col]
    end
    
    def cells_by_row(row)
        @cells.select { |cell| cell.row == row }
    end

    def cells_by_col(col)
        @cells.select { |cell| cell.col == col }
    end

    def cells_by_box(box)
        @cells.select { |cell| cell.box == box }
    end
    
    def cell_by_ref(ref)
        @cells.select { |cell| cell.ref == ref }[0]
    end
    
    def singles
        found = false
        
        (0..8).each do |row|
            "123456789".each_char do |possible|
                cells = cells_by_row(row).select { |cell| cell.has?(possible) }
                if cells.length == 1
                    puts "Single #{possible} in row #{row}"
                    solve(cells[0], possible)
                    found = true
                end
            end
        end
        
        (0..8).each do |col|
            "123456789".each_char do |possible|
                cells = cells_by_col(col).select { |cell| cell.has?(possible) }
                if cells.length == 1
                    puts "Single #{possible} in col #{col}"
                    solve(cells[0], possible)
                    found = true
                end
            end
        end
        
        (0..8).each do |box|
            "123456789".each_char do |possible|
                cells = cells_by_box(box).select { |cell| cell.has?(possible) }
                if cells.length == 1
                    puts "Single #{possible} in box #{box}"
                    solve(cells[0], possible)
                    found = true
                end
            end
        end
        
        found
    end

    def naked_cells( cells )
        found = false
        
        cells = cells.select { |cell| cell.possibles.length > 1 }
        possibles = cells.map { |cell| cell.possibles }.join.split(//).uniq.sort
        
        combos = []
        (2...possibles.length).each do |r|
            possibles.combination(r) { |x| combos << x }
        end
        
        combos.each do |combo|
            matches = cells.select { |cell| (cell.possibles.split(//) - combo).length == 0 }
            if matches.length == combo.length
                remove = combo.join
                cells.each do |cell|
                    if not matches.include? cell
                        found = true if cell.remove_possibles(remove)
                    end
                end
                if found
                    where = matches.map { |cell| "#{cell.row}:#{cell.col}" }.join(" ")
                    puts "Naked #{remove} found in #{where}"
                    return true
                end
            end
        end
        
        found
    end
    
    def nakeds
        # http://www.sudokuwiki.org/Naked_Candidates#NP
        (0..8).each { |row| return true if naked_cells(cells_by_row(row)) }
        (0..8).each { |col| return true if naked_cells(cells_by_col(col)) }
        (0..8).each { |box| return true if naked_cells(cells_by_box(box)) }
        false
    end
    
    def pointing_pairs
        # http://www.sudokuwiki.org/Intersection_Removal#IR
        found = false
        row = col = 0
        
        (0..8).each do |box|
            "123456789".each_char do |possible|
                pairs = cells_by_box(box).select { |cell| cell.has? possible}
                if pairs.length == 2 or pairs.length == 3
                    if pairs[0].row == pairs[1].row and pairs[0].row == pairs[-1].row
                        cells_by_row(pairs[0].row).each do |cell| 
                            if not pairs.include? cell and cell.remove(possible) 
                                found = true
                            end
                        end
                    end
                    if pairs[0].col == pairs[1].col and pairs[0].col == pairs[-1].col
                        cells_by_col(pairs[0].col).each do |cell| 
                            if not pairs.include? cell and cell.remove(possible) 
                                found = true
                            end
                        end
                    end
                    if found
                        puts "Pointing Pair #{possible} on box #{box}"
                        return true
                    end
                end
            end
        end
        
        found        
    end
        
    def box_line_reduction
        # http://www.sudokuwiki.org/Intersection_Removal#LBR
        found = false
        row = col = 0
        
        (0..8).each do |row|
            "123456789".each_char do |possible|
                cells = cells_by_row(row).select { |cell| cell.has? possible }
                boxes = cells.map { |cell| cell.box}.uniq
                if boxes.length == 1
                    cells_by_box(boxes[0]).each do |cell| 
                        if not cells.include? cell and cell.remove(possible) 
                            found = true
                        end
                    end
                end
                if found
                    puts "Box Line reduction of #{possible} on row #{row}"
                    return true
                end
            end
        end
        
        (0..8).each do |col|
            "123456789".each_char do |possible|
                cells = cells_by_col(col).select { |cell| cell.has? possible }
                boxes = cells.map { |cell| cell.box}.uniq
                if boxes.length == 1
                    cells_by_box(boxes[0]).each do |cell| 
                        if not cells.include? cell and cell.remove(possible) 
                            found = true
                        end
                    end
                end
                if found
                    puts "Box Line reduction of #{possible} on col #{col}"
                    return true
                end
            end
        end
        
        found        
    end
        
    def swordfishes( scale )
        # http://www.sudokuwiki.org/Sword_Fish_Strategy

        # Extension of swordfish to allow variations apart from 3x3
        # Note: 2x2 should be same as X-Wing

        type = "X-Wing" if scale == 2
        type = "Swordfish" if scale == 3
        type = "Swordfish-#{scale}" if scale > 3

        found = false
        matches = []

        "123456789".each_char do |possible|

            # Find rows with this possible. 
            # Store in a tuple [ possible, row, columns … ]

            rows = []        
            (0..8).each do |row|
                cols = [possible, row] + cells_by_row(row).select { |cell| cell.has? possible }.map { |cell| cell.col }
                rows << cols if cols.length >= 3 && cols.length <= (scale+2)
            end
            
            # If we have at least {scale} rows then try to find a combination of scale rows that share the same {scale} columns
            # First get all combinations of (scale) rows that contain our possible

            combos = []
            if rows.size >= scale
                rows.map {|row| row[1]}.combination(scale) { |x| combos << x }
            end
            
            # Now check if the (scale) rows share the same (scale) cols
            combos.each do |combo|
                scale_rows = rows.select {|row| combo.include? row[1] }
                cols = []
                scale_rows.each {|t| cols += t[2..-1]}
                cols.uniq!

                if scale_rows.size == scale && cols.size == scale
                    # Remove possibles from all the other rows
                    row_numbers = scale_rows.map {|t| t[1] }
                    cols.each do |col|
                        cells_by_col(col).each do |cell|
                            if not row_numbers.include? cell.row
                                found = true if cell.remove(possible)
                            end
                        end
                    end
                    if found
                        puts "#{type} #{possible} in rows #{row_numbers.inspect} cols #{cols.inspect}"
                        return true
                    end
                end
            end
        end
        
        "123456789".each_char do |possible|

            # Find cols with this possible. 
            # Store in a tuple [ possible, col, rows … ]


            cols = []        
            (0..8).each do |col|
                rows = [possible, col] + cells_by_col(col).select { |cell| cell.has? possible }.map { |cell| cell.row }
                cols << rows if rows.length >= 3 && rows.length <= (scale+2)
            end
            
            # If we have at least {scale} cols then try to find a combination of {scale} cols that share the same {scale} rows
            # First get all combinations of {scale} cols that contain our possible

            combos = []
            if cols.size >= scale
                cols.map {|col| col[1]}.combination(4) { |x| combos << x }
            end
            
            # Now check if the {scale} cols share the same {scale} rows
            combos.each do |combo|
                scale_cols = cols.select {|col| combo.include? col[1] }
                rows = []
                scale_cols.each {|t| rows += t[2..-1]}
                rows.uniq!
                if scale_cols.size == scale && rows.size == scale
                    # Remove possibles from all the other cols
                    col_numbers = scale_cols.map {|t| t[1] }
                    rows.each do |row|
                        cells_by_row(row).each do |cell|
                            if not col_numbers.include? cell.col
                                found = true if cell.remove(possible)
                            end
                        end
                    end
                    if found
                        puts "#{type} #{possible} in cols #{col_numbers.inspect} rows #{rows.inspect}"
                        return true
                    end
                end
            end
            
        end
        
        found        
    end
    
    def solved
        @cells.all? { |cell| cell.possibles.length == 1 }
    end

    def solution
        if solved
            return @cells.map { |cell| cell.possibles }.join
        else
            return "Not solved"
        end
    end
        
end

line = 0
solved=0
total=0
solutions = []
solution = ""


File.new("top95expected.txt", "r").each do |data|
   solutions << solution
end

solutions = File.readlines("top95expected.txt")

line = 0
File.new("top95.txt", "r").each do |data|
    exit if data.length < 81
    line += 1
    puts "=== Populating board #{line} ==="
    board = Board.new(data)
    running = true
    while running
        if board.singles
        elsif board.nakeds
        elsif board.pointing_pairs
        elsif board.box_line_reduction
        elsif board.swordfishes(2)
        elsif board.swordfishes(3)
        elsif board.swordfishes(4)
        else 
            running = false
        end
    end
    puts board.to_s
    if board.solved
        solved += 1
        expected = solutions[line-1].strip
        if board.solution != expected
            puts "Houston we have a problem … expected\n[#{expected}] \ngot\n[#{board.solution}]"
            exit
        end
        puts ">>> SOLVED >>>"
    else
        puts ">>> BEATEN >>>"
    end
    total += 1
end

puts ">>> Solved #{solved} of #{total} puzzles"



