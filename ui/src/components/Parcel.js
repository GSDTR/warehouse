/* @flow */

//*
import React, {Component} from 'react';
import "./Parcel.css"

var invariant = require('invariant');

var IDX_OF_ISSUE_IDX = 0;
var IDX_ESTIMATION_TIME = 4;

const body = document.body;

type Props = {
    search: boolean,
    headers: Array<string>,
    initialData: Array<Array<string>>,
    columnClasses: Array<string>,
    columnDataTypes: Array<string>,
};

type EditState = {
    row: number,
    cell: number,
};

type State = {
    data: Array<Array<string>>,
    sortby: ?string,
    descending: boolean,
    edit: ?EditState, // [row index, cell index],
    search: boolean,
    hideCompleted: boolean,
};

class Parcel extends Component<Props, State>{
    displayName: 'Excel';
    _preSearchData: Array<Array<string>>;
    _searchFields: Array<string>;
    _sortingColumn: number;
    _resultedRow: Array<string>;
    rowIdxToDisplay: number;

    _focusElem: ?HTMLInputElement;

    constructor(props: Props) {
        super(props);
            this._preSearchData = this.props.initialData;
        this.state = {
            data: this.props.initialData,
            sortby: null,
            descending: false,
            edit: null, // [row index, cell index],
            search: this.props.search,
            hideCompleted: false,
        };
        this._searchFields = Array(this.props.headers.length).join(".").split(".");
        this._resultedRow = Array(this.props.headers.length).join(".").split(".");
//        console.log("1. resulted row: ", this._resultedRow);
        (this:any)._search = this._search.bind(this);
        (this:any).search_loop = this.search_loop.bind(this);
        (this:any)._showEditor = this._showEditor.bind(this);
        (this:any)._renderSearch = this._renderSearch.bind(this);
        (this:any)._sort = this._sort.bind(this);
        (this:any)._save = this._save.bind(this);
        (this:any).requestParcelsList = this.requestParcelsList.bind(this);
    }

    hideCompleted(hideCompletedFlag: boolean) {
        this.setState( {
            hideCompleted: hideCompletedFlag
        });
    }

    componentWillMount() {
//        console.log("componentWillMount");
    }

    componentDidMount() {
//        console.log("componentDidMount");
//        document.title = "Warehouse";
        setInterval(this.requestParcelsList, 3000);
//        this.fetchData("/api/v1/parcels");
    }

    requestParcelsList() {
//        console.log("requestParcelsList");
        this.fetchData("/api/v1/parcels");
    }

    updateData(newData: Array<Array<string>>) {
        var dataJson = newData;
        for(let i = 0; i < newData.length; i++){
            let childArray = newData[i];
            for(let j = 0; j < childArray.length; j++){
                if (childArray[j] == null) {
                    console.log("null", childArray[j]);
                    childArray[j] = "";
                }
            }
        }
        // console.log("dataE: ", newData);
        this.setState({
                data: dataJson
            });
        this._preSearchData = dataJson;
        var tmpData = this.search_loop(this._preSearchData);
        this.sort_data(tmpData, this._sortingColumn, false);
    }

    sort_data(data, column, invert = true) {
        this._sortingColumn = column;
//        var data = this.state.data.slice();
//        console.log("data to sort len: ", data.length);
        var descending;
        if (invert) {
            descending = this.state.sortby === column && !this.state.descending;
        } else {
            descending = this.state.sortby === column && this.state.descending;
        }
            data.sort(  (a, b) => {
                var aa = 0;
                var bb = 0;
                var col_name =  this.props.headers[column];
                // console.log("column: ", col_name, this.props.columnDataTypes[col_name] );
                if ((this.props.columnDataTypes[column] == 'numeric') || (this.props.columnDataTypes[column] == 'int')){
                    // console.log("a: ", a[col_name], "b: ", b[col_name]);
                    if (a[col_name] != "") {
                        aa = parseInt(a[col_name], 10);
                    }
                    if (b[col_name] != "") {
                        bb = parseInt(b[col_name], 10);
                    }
                } else if (this.props.columnDataTypes[column] == 'alphaNumeric') {

                    var reSingleDigit = /(?<!\d)[\d](?!\d)/g;
                    //console.log("a: ", a[col_name], "b: ", b[col_name]);
                    aa = a[col_name].replace(reSingleDigit, "0$&");
                    bb = b[col_name].replace(reSingleDigit, "0$&");
                    //console.log(a[col_name], " --> ", aa, ". ", b[col_name], " --> ", bb);

                } else {
                    // console.log("col: ", col_name, "a: ", a[ this.props.headers[col_name] ], "b: ", b[ this.props.headers[col_name] ]);
                    aa = a[col_name];
                    bb = b[col_name];
                }
                if (descending == false) {
                    if ( aa < bb ) {    return -1;  }
                    if ( aa > bb ) {    return 1;   }
                    return 0;
                } else {
                    if ( aa < bb ) {    return 1;   }
                    if ( aa > bb ) {    return -1;  }
                    return 0;
                }
            });
        this.setState({
            data: data,
            sortby: column,
            descending: descending,
        });
    }

    _sort(e: any) {
        var column = e.target.cellIndex;
        var data = this.state.data.slice();
        this.sort_data(data, column);

    }

    _showEditor(e: any) {
        this.setState({edit: {
                row: parseInt(e.target.dataset.row, 10),
                cell: e.target.cellIndex,
            }});
    }

    postData(url, data) {
        var request = new XMLHttpRequest();
        request.open('POST', url, true);
        request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

        request.onreadystatechange = function() {//Вызывает функцию при смене состояния.
           if(request.readyState == XMLHttpRequest.DONE && request.status == 200) {
                // Запрос завершен. Здесь можно обрабатывать результат.
            }
}
        request.send(JSON.stringify(data));

    }

    fetchData(url) {
        var xhr = new XMLHttpRequest();
        xhr.responseType = 'json';
        xhr.onload = () => {
            var resp = xhr.response;
            if(resp["success"] != "true") {
                alert(resp["error"]);
            } else {
                if( "data" in resp) {
                    var dataa = JSON.parse(resp["data"]);
                    if(dataa != null) {
                        this.updateData(dataa);
                    } else {
                    }
                } else {
                }
            }
        };
        xhr.open('GET', url, true);
        xhr.send(null);

    }


    _save(e: any) {
        e.preventDefault();
        var input = e.target.firstChild;
        var data = this.state.data.slice();
        var rowIdx = this.state.edit.row;
        var colIdx = this.state.edit.cell;
        var colName = this.props.headers[colIdx];
        var val = input.value;
        var id = this.state.data[rowIdx][ this.props.headers[0] ];
        var part_number = this.state.data[rowIdx][ this.props.headers[4] ];
        console.log("param ", colName, " val ", val);
        this.fetchData("/api/v1/cell/"+id+"?action=update&" + "param=" + colName + "&val="+encodeURIComponent(val)+"&part_number="+part_number);



        invariant(this.state.edit, "edit field can't be undefined here!");
        data[this.state.edit.row][this.state.edit.cell] = input.value;
        this.setState({
            edit: null,
            data: data,
        });
    }


    search_loop(tmpData) {
        var idx = 0;
        var needle = "";
        // console.log("tmpData: ", tmpData);
        for (var i = 0; i < this._searchFields.length; i++) {
            idx = i;
            needle = this._searchFields[i];
            if (needle == "") {
                continue;
            }
            tmpData = tmpData.filter( (row) => {
                if ((this.props.columnDataTypes[i] == 'numeric') || (this.props.columnDataTypes[i] == 'int')){
//                    console.log("numeric search");
                    return row[this.props.headers[idx]] == needle;
                } else {
//                    console.log("string search");
                    return row[this.props.headers[idx]].toString().toLowerCase().indexOf(needle) > -1;
                }
//                return row[ this.props.headers[idx] ].toString().toLowerCase().indexOf(needle) > -1;
            });
        }
//        console.log("search res dataset length: ", tmpData.length);
        return tmpData;
    }

    _search(e: any) {
        var needle = e.target.value.toLowerCase();
        var idx = e.target.dataset.idx;
        this._searchFields[idx] = needle;

//        console.log("search: ", this._searchFields);
        var tmpData = this._preSearchData;

        tmpData = this.search_loop(tmpData);
        this.sort_data(tmpData, this._sortingColumn, false);
    }


    exportBlob(format: string, contents: string) {
        var URL = window.URL || window.webkitURL;
        var blob = new Blob([contents], {type: 'text/' + format});

        var a = document.createElement("a"),
            url = URL.createObjectURL(blob);
        a.href = url;
        a.download = "data." + format;
        if(body != null) {
            body.appendChild(a);
        }
        a.click();
        setTimeout(function() {
            if(body != null) {
                body.removeChild(a);
            }
            window.URL.revokeObjectURL(url);
        }, 0);

    }

    exportJSON() {
        var format = "json";
        var contents = JSON.stringify(this.state.data);
        this.exportBlob(format, contents);
    }

    exportCSV() {
        var format = "csv";
        var contents = this.state.data.reduce(function(result, row) {
            return result
                + row.reduce(function(rowresult, cell, idx) {
                    return rowresult
                        + '"'
                        + cell.replace(/"/g, '""')
                        + '"'
                        + (idx < row.length - 1 ? ',' : '');
                }, '')
                + "\n";
        }, '');
        this.exportBlob(format, contents);
    }

    render() {
        return (
            <div className="Excel">
                {this._renderTable()}
            </div>
        );
    }

    _renderSearch() {
        if (!this.state.search) {
            return null;
        }
        return (
            <tr onChange={this._search}>
                {this.props.headers.map((_ignore, idx) => {
                    return <td key={idx} class={this.props.columnClasses[idx]}><input type="text" data-idx={idx}/></td>;
                })}
            </tr>
        );
    }

    _renderResultRow() {
//        console.log("4. resulted row: ", this._resultedRow);
        return (
            <div>
            <table>
            <tr onChange={this._search}>
                {this.props.headers.map((_ignore, idx) => {
                    if( idx == this.props.headers.length-1) {
                        return <td key={idx} class={this.props.columnClasses[idx] + " resRow"}>{this.rowIdxToDisplay}</td>;
                    } else {
                        return <td key={idx}
                                   class={this.props.columnClasses[idx] + " resRow"}>{this._resultedRow[idx]}</td>;
                    }
                })}
            </tr>
            </table>
            </div>
        );
    }

    componentDidUpdate(prevProps, prevState) {
//        console.log("Did update");
//        console.log("focus elem", this._focusElem);
        if( this._focusElem) {
            this._focusElem.focus();
        }
    }

    _renderTable() {
        this.rowIdxToDisplay = 0;
        var i;
        for(i=0; i<this._resultedRow.length; i++) {
//            console.log("col ", i, this.props.columnDataTypes[i]);
            if(this.props.columnDataTypes[i] == "int") {
                this._resultedRow[i] = 0;
            } else if(this.props.columnDataTypes[i] == "serial") {
                this._resultedRow[i] = 0;
            }
        }
//        console.log("3. resulted row: ", this._resultedRow);

        return (
            <div>
            <div>
            <table>
                <thead onClick={this._sort}>
                <tr>{
                    this.props.headers.map((title, idx) => {
//                        console.log("idx", idx, title);
                        if (this.state.sortby === idx) {
                            title += this.state.descending ? ' \u2191' : ' \u2193';
                        }
                        return <th key={idx} class={this.props.columnClasses[idx]}>{title}</th>;
                    }, this)
                }</tr>
                </thead>
            </table>
            </div>
            <div  class="tbl">
            <table>
                <tbody onDoubleClick={this._showEditor}>
                {this._renderSearch()}
                {this.state.data.map(function (row, rowidx) {
                    // console.log("row: ", row);
                    var rowClass = "notFinished";
                    this.rowIdxToDisplay += 1;
                    return (
                        <tr key={rowidx} class={rowClass}>{
//                            Object.values(row).map((cell, idx) => {
                            this.props.headers.map((cell, idx) => {
                                var content = this.state.data[rowidx][cell];
                                var edit = this.state.edit;
                                if (edit && edit.row === rowidx && edit.cell === idx) {
                                    content = (
                                        <form onSubmit={this._save}>
                                            <input type="text" defaultValue={content} ref={c => (this._focusElem = c)} />
                                        </form>
                                    );
                                }
                                 if( this.props.columnDataTypes[idx] == "int" ) {
                                    this._resultedRow[idx] += parseInt(content) || 0;
                                } else if(this.props.columnDataTypes[idx] == "serial") {
                                    this._resultedRow[idx] = this.rowIdxToDisplay;
                                }
                                if (this.props.columnDataTypes[idx] == 'link') {
                                    if (edit && edit.row === rowidx && edit.cell === idx) {
                                        return <td class={this.props.columnClasses[idx]} key={idx} data-row={rowidx}>{content}</td>;
                                    }
                                    var myRegexp = /^(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+)/img;
                                    var match = myRegexp.exec(content);
                                    var title = " ";
                                    if (match != null) {
                                        title = match[1];
                                        //console.log("idx", idx, content, title);
                                        return <td class={this.props.columnClasses[idx]} key={idx} data-row={rowidx}><a href={content} class={rowClass} target="_blank">{title}</a></td>;
                                    } else {
                                        //console.log("idx: ", idx, "content: ", content, "title: ", title);
                                        return <td class={this.props.columnClasses[idx]} key={idx} data-row={rowidx}>{content}</td>;
                                    }
                                } else if (this.props.columnDataTypes[idx] == 'serial') {
                                    return <td class={this.props.columnClasses[idx]} key={idx} data-row={rowidx}>{this.rowIdxToDisplay}</td>;
                                } else {
                                    return <td class={this.props.columnClasses[idx]} key={idx} data-row={rowidx}>{content}</td>;
                                }
                            }, this)}
                        </tr>
                    );
                }, this)}
                </tbody>
            </table>
            </div>
            {this._renderResultRow()}
            </div>
        );
    }
}

export default Parcel


